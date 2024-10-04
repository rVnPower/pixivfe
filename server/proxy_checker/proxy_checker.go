package proxy_checker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"codeberg.org/vnpower/pixivfe/v2/config"
	"codeberg.org/vnpower/pixivfe/v2/server/utils"
)

const (
	testImagePath = "/img-original/img/2024/01/21/20/50/51/115365120_p0.jpg"
)

var (
	workingProxies      []string
	workingProxiesMutex sync.RWMutex
)

func CheckProxies(ctx context.Context) {
	logln("Starting proxy check...")
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var newWorkingProxies []string

	logf("Total proxies to check: %d", len(config.BuiltinProxyList))

	for _, proxy := range config.BuiltinProxyList {
		wg.Add(1)
		go func(proxyURL string) {
			defer wg.Done()
			isWorking, resp := testProxy(ctx, proxyURL)
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

func testProxy(ctx context.Context, proxyBaseURL string) (bool, *http.Response) {
	fullURL := fmt.Sprintf("%s%s", strings.TrimRight(proxyBaseURL, "/"), testImagePath)
	logf("Testing proxy %s with full URL: %s", proxyBaseURL, fullURL)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
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
