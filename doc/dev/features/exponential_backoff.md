# Exponential Backoff in PixivFE

PixivFE implements exponential backoff as a strategy for handling API request failures and managing token usage. This document outlines the implementation details of this feature across two key components of the application.

## Overview

Exponential backoff is a technique used to gradually increase the wait time between retries of a failed operation. In PixivFE, this technique is applied at two levels:

1. API request level
2. Token management level

### Key aspects

1. Both implementations use a base timeout value that doubles with each retry or failure.
2. Both have a maximum backoff time to prevent excessively long waits.
3. The token manager's implementation is per-token, allowing for fine-grained control over individual token usage.
4. The request retry mechanism applies to all requests, providing a general-purpose backoff strategy.
5. Backoff settings for both levels can be configured separately via environment variables.

## Configuration

The backoff settings for both API request level and token management level can be configured using environment variables. If not set, default values are used.

### API request level backoff

- `PIXIVFE_API_MAX_RETRIES`: Maximum number of retries for API requests
- `PIXIVFE_API_BASE_TIMEOUT`: Base timeout duration for API requests
- `PIXIVFE_API_MAX_BACKOFF_TIME`: Maximum backoff time for API requests

### Token management level backoff

- `PIXIVFE_MAX_RETRIES`: Maximum number of retries for token management
- `PIXIVFE_BASE_TIMEOUT`: Base timeout duration for token management
- `PIXIVFE_MAX_BACKOFF_TIME`: Maximum backoff time for token management

## API request level backoff

**Location: `core/requests.go`**

The `retryRequest` function implements exponential backoff for API requests:

- It attempts to make a request up to `config.GlobalConfig.APIMaxRetries` times.
- If a request fails, it calculates the backoff duration using the formula:
  ```go
  backoffDuration := time.Duration(float64(config.GlobalConfig.APIBaseTimeout) * float64(1<<uint(i)))
  ```
  This effectively doubles the backoff duration with each retry attempt.
- The backoff duration is capped at a maximum value:
  ```go
  if backoffDuration > config.GlobalConfig.APIMaxBackoffTime {
      backoffDuration = config.GlobalConfig.APIMaxBackoffTime
  }
  ```
- The function then waits for the calculated backoff duration before making the next attempt.

This approach helps to avoid overwhelming the Pixiv API with rapid subsequent requests after a failure, giving it time to recover.

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
