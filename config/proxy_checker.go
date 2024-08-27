package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	proxyCheckTimeout = 10 * time.Second
	testImagePath     = "/img-original/img/2024/01/21/20/50/51/115365120_p0.jpg"
)

var (
	workingProxies      []string
	workingProxiesMutex sync.RWMutex
	stopChan            chan struct{}
)

func InitializeProxyChecker(r context.Context) {
	stopChan = make(chan struct{})
	StartProxyChecker(r)
}

func StartProxyChecker(r context.Context) {
	go func() {
		for {
			select {
			case <-stopChan:
				log.Println("Stopping proxy checker...")
				return
			default:
				checkProxies(r)
				if t := GlobalServerConfig.ProxyCheckInterval; t > 0 {
					time.Sleep(t)
				} else {
					log.Println("Proxy check interval set to 0, disabling auto-check from now on.")
					select {} // Sweet dreams!
				}
			}
		}
	}()
}

func StopProxyChecker() {
	close(stopChan)
}

func checkProxies(r context.Context) {
	logln("Starting proxy check...")
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var newWorkingProxies []string

	logf("Total proxies to check: %d", len(BuiltinProxyList))

	for _, proxy := range BuiltinProxyList {
		wg.Add(1)
		go func(proxyURL string) {
			defer wg.Done()
			isWorking, resp := testProxy(r, proxyURL)
			status := ""
			if resp != nil {
				status = resp.Status
			}
			if isWorking {
				mutex.Lock()
				newWorkingProxies = append(newWorkingProxies, proxyURL)
				mutex.Unlock()
				logf("[OK]  %s %s", proxyURL, status)
			} else {
				logf("[ERR] %s %s", proxyURL, status)
			}
		}(proxy)
	}

	wg.Wait()

	updateWorkingProxies(newWorkingProxies)
}

func testProxy(r context.Context, proxyBaseURL string) (bool, *http.Response) {
	client := &http.Client{Timeout: proxyCheckTimeout}

	fullURL := fmt.Sprintf("%s%s", strings.TrimRight(proxyBaseURL, "/"), testImagePath)
	logf("Testing proxy %s with full URL: %s", proxyBaseURL, fullURL)

	req, err := http.NewRequestWithContext(r, "GET", fullURL, nil)
	if err != nil {
		logf("Error creating request for proxy %s: %v", proxyBaseURL, err)
		return false, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		logf("Error testing proxy %s: %v", proxyBaseURL, err)
		return false, nil
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, resp
}

func updateWorkingProxies(newProxies []string) {
	workingProxiesMutex.Lock()
	defer workingProxiesMutex.Unlock()

	workingProxies = newProxies
	logf("Updated working proxies. Count: %d", len(workingProxies))
}

func GetWorkingProxies() []string {
	workingProxiesMutex.RLock()
	defer workingProxiesMutex.RUnlock()

	return append([]string{}, workingProxies...)
}

// Helper functions for logging
func logf(format string, v ...any) {
	log.Printf(format, v...)
}

func logln(v ...any) {
	log.Println(v...)
}
