package proxy_checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
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

	isWorking, resp := testProxy(context.Background(), server.URL)
	if !isWorking {
		t.Errorf("Expected proxy to be working")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	isWorking, resp = testProxy(context.Background(), "http://nonexistentproxy.invalid")
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

	isWorking, resp = testProxy(context.Background(), errorServer.URL)
	if isWorking {
		t.Errorf("Expected proxy to be not working due to error status")
	}
	if resp == nil || resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected error status response")
	}
}
