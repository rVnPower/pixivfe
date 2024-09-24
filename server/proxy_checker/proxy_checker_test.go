package proxy_checker

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"codeberg.org/vnpower/pixivfe/v2/config"
)

func TestUpdateAndGetWorkingProxies(t *testing.T) {
	t.Parallel()
	newProxies := []string{"http://proxy1.invalid", "http://proxy2.invalid"}
	updateWorkingProxies(newProxies)

	result := GetWorkingProxies()
	if len(result) != len(newProxies) {
		t.Errorf("Expected %d proxies, got %d", len(newProxies), len(result))
	}

	for i, proxy := range result {
		if proxy != newProxies[i] {
			t.Errorf("Expected proxy %s, got %s", newProxies[i], proxy)
		}
	}

	// Test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = GetWorkingProxies()
		}()
	}
	wg.Wait()
}

func TestTestProxy(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	isWorking, resp := testProxy(server.URL)
	if !isWorking {
		t.Errorf("Expected proxy to be working")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	isWorking, resp = testProxy("http://nonexistentproxy.invalid")
	if isWorking {
		t.Errorf("Expected proxy to be not working")
	}
	if resp != nil {
		t.Errorf("Expected nil response for non-working proxy")
	}

	// Test with a server that returns an error status
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	isWorking, resp = testProxy(errorServer.URL)
	if isWorking {
		t.Errorf("Expected proxy to be not working due to error status")
	}
	if resp == nil || resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected error status response")
	}
}

func TestProxyCheckerBehavior(t *testing.T) {
	// Set up test configuration
	config.GlobalConfig.ProxyCheckInterval = 10 * time.Millisecond
	config.BuiltinProxyList = []string{"http://proxy1.invalid", "http://proxy2.invalid"}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		done := InitializeProxyChecker()
		<-done // Wait for first check to complete

		time.Sleep(15 * time.Millisecond) // Wait for one more check

		StopProxyChecker()
	}()

	wg.Wait()

	proxies := GetWorkingProxies()
	if len(proxies) > 0 {
		t.Errorf("Expected no working proxies, got %d", len(proxies))
	}
}

func TestInitializeProxyChecker(t *testing.T) {
	config.GlobalConfig.ProxyCheckInterval = 10 * time.Millisecond
	config.BuiltinProxyList = []string{"http://proxy1.invalid", "http://proxy2.invalid"}

	done := InitializeProxyChecker()
	<-done // Wait for first check to complete

	time.Sleep(15 * time.Millisecond) // Wait for one more check

	StopProxyChecker()

	// Ensure that the proxy checker can be restarted
	done = InitializeProxyChecker()
	<-done

	StopProxyChecker()
}

func TestProxyCheckerWithZeroInterval(t *testing.T) {
	config.GlobalConfig.ProxyCheckInterval = 0
	config.BuiltinProxyList = []string{"http://proxy1.invalid", "http://proxy2.invalid"}

	done := InitializeProxyChecker()
	<-done // Wait for first check to complete

	initialProxies := len(GetWorkingProxies())

	time.Sleep(50 * time.Millisecond) // Wait to ensure no more checks are performed

	if len(GetWorkingProxies()) != initialProxies {
		t.Errorf("Expected proxy count to remain %d, but got %d", initialProxies, len(GetWorkingProxies()))
	}

	// StopProxyChecker should not cause any issues
	StopProxyChecker()
}
