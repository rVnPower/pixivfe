package proxy_checker

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

const (
	testImagePath = "/img-original/img/2024/01/21/20/50/51/115365120_p0.jpg"
)

var (
	workingProxies      []string
	workingProxiesMutex sync.RWMutex
	stopChan            chan struct{} = make(chan struct{})
)

func InitializeProxyChecker() chan struct{} {
	firstCheckDone := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopChan:
				log.Print("Stopping proxy checker...")
				return
			default:
				checkProxies()
				select {
				case <-firstCheckDone:
					// First check already done, do nothing
				default:
					close(firstCheckDone) // Signal that the first check is done
				}
				if t := config.GlobalConfig.ProxyCheckInterval; t > 0 {
					time.Sleep(t)
				} else {
					log.Print("Proxy check interval set to 0, disabling auto-check from now on.")
					select {} // Sweet dreams!
				}
			}
		}
	}()
	return firstCheckDone
}

func StopProxyChecker() {
	close(stopChan)
}

func checkProxies() {
	logln("Starting proxy check...")
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var newWorkingProxies []string

	logf("Total proxies to check: %d", len(config.BuiltinProxyList))

	for _, proxy := range config.BuiltinProxyList {
		wg.Add(1)
		go func(proxyURL string) {
			defer wg.Done()
			isWorking, resp := testProxy(proxyURL)
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

func testProxy(proxyBaseURL string) (bool, *http.Response) {
	fullURL := fmt.Sprintf("%s%s", strings.TrimRight(proxyBaseURL, "/"), testImagePath)
	logf("Testing proxy %s with full URL: %s", proxyBaseURL, fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		logf("Error creating request for proxy %s: %v", proxyBaseURL, err)
		return false, nil
	}

	resp, err := utils.HttpClient.Do(req)
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
	log.Print(v...)
}
