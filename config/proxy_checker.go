package config

import (
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

func InitializeProxyChecker() {
	stopChan = make(chan struct{})
	StartProxyChecker()
}

func StartProxyChecker() {
	go func() {
		for {
			select {
			case <-stopChan:
				log.Println("Stopping proxy checker...")
				return
			default:
				checkProxies()
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

func checkProxies() {
	logln("Starting proxy check...")
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var newWorkingProxies []string

	logf("Total proxies to check: %d", len(BuiltinProxyList))

	for _, proxy := range BuiltinProxyList {
		wg.Add(1)
		go func(proxyURL string) {
			defer wg.Done()
			if isProxyWorking(proxyURL) {
				mutex.Lock()
				newWorkingProxies = append(newWorkingProxies, proxyURL)
				mutex.Unlock()
				logf("Working proxy found: %s", proxyURL)
			} else {
				logf("Non-working proxy: %s", proxyURL)
			}
		}(proxy)
	}

	wg.Wait()

	updateWorkingProxies(newWorkingProxies)
}

func isProxyWorking(proxyBaseURL string) bool {
	client := &http.Client{Timeout: proxyCheckTimeout}

	fullURL := fmt.Sprintf("%s%s", strings.TrimRight(proxyBaseURL, "/"), testImagePath)
	logf("Testing proxy %s with full URL: %s", proxyBaseURL, fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		logf("Error creating request for proxy %s: %v", proxyBaseURL, err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		logf("Error testing proxy %s: %v", proxyBaseURL, err)
		return false
	}
	defer resp.Body.Close()

	logf("Proxy %s response status: %s", proxyBaseURL, resp.Status)

	return resp.StatusCode == http.StatusOK
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
func logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func logln(v ...interface{}) {
	log.Println(v...)
}
