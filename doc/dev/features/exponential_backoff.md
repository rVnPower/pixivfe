# Exponential Backoff in PixivFE

PixivFE implements exponential backoff as a strategy for handling API request failures and managing token usage. This document outlines the implementation details of this feature across two key components of the application.

## Overview

Exponential backoff is a technique used to gradually increase the wait time between retries of a failed operation. In PixivFE, this technique is applied at two levels:

1. API request level
2. Token management level

## API request level backoff

**Location: `core/requests.go`**

PixivFE uses the `retryablehttp` package to implement exponential backoff for API requests. The implementation is as follows:

1. A `retryablehttp.Client` is initialized in the `init` function:

```go
func init() {
    retryClient = retryablehttp.NewClient()
    retryClient.RetryMax = config.GlobalConfig.APIMaxRetries
    retryClient.RetryWaitMin = config.GlobalConfig.APIBaseTimeout
    retryClient.RetryWaitMax = config.GlobalConfig.APIMaxBackoffTime
    retryClient.HTTPClient = utils.HttpClient
}
```

2. The `retryRequest` function uses this client to perform requests with automatic retries:

```go
func retryRequest(ctx context.Context, reqFunc func(context.Context, string) (*retryablehttp.Request, error)) (SimpleHTTPResponse, error) {
    // ... (function implementation)
}
```

## Token management level backoff

**Location: `server/token_manager/token_manager.go`**

The `TokenManager` implements exponential backoff for individual tokens:

- In the `MarkTokenStatus` method, when a token is marked as `TimedOut`, it calculates a timeout duration:
  ```go
  timeoutDuration := time.Duration(math.Min(
      float64(tm.baseTimeout)*math.Pow(2, float64(token.FailureCount-1)),
      float64(tm.maxBackoffTime),
  ))
  ```
- This calculation uses the `math.Pow` function to implement exponential growth based on the number of consecutive failures.
- The timeout duration is also capped at a maximum value (`tm.maxBackoffTime`).
- The token's `TimeoutUntil` is set to the current time plus this calculated duration.

This approach allows tokens that repeatedly fail increasingly longer "cool-down" periods before being used again, helping to manage rate limiting of individual tokens by the Pixiv API.

## Implementation details

### Configuration (`config/config.go`)

The `ServerConfig` struct in `config/config.go` includes fields for both API request level and token management level backoff settings:

```go
type ServerConfig struct {
    // ... other fields ...
    MaxRetries     int           `env:"PIXIVFE_MAX_RETRIES,overwrite"`
    BaseTimeout    time.Duration `env:"PIXIVFE_BASE_TIMEOUT,overwrite"`
    MaxBackoffTime time.Duration `env:"PIXIVFE_MAX_BACKOFF_TIME,overwrite"`

    APIMaxRetries     int           `env:"PIXIVFE_API_MAX_RETRIES,overwrite"`
    APIBaseTimeout    time.Duration `env:"PIXIVFE_API_BASE_TIMEOUT,overwrite"`
    APIMaxBackoffTime time.Duration `env:"PIXIVFE_API_MAX_BACKOFF_TIME,overwrite"`
    // ... other fields ...
}
```

The `LoadConfig` method sets default values for these settings if they are not provided through environment variables:

```go
func (s *ServerConfig) LoadConfig() error {
    // ... other initializations ...
    s.MaxRetries = 5
    s.BaseTimeout = 1 * time.Second
    s.MaxBackoffTime = 32 * time.Second

    s.APIMaxRetries = 3
    s.APIBaseTimeout = 500 * time.Millisecond
    s.APIMaxBackoffTime = 8 * time.Second
    // ... rest of the method ...
}
```
