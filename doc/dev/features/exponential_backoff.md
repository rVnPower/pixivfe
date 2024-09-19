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

## API request level backoff

**Location: `core/requests.go`**

The `retryRequest` function implements exponential backoff for API requests:

- It attempts to make a request up to `tokenManager.GetMaxRetries()` times.
- If a request fails, it calculates the backoff duration using the formula:
  ```go
  backoffDuration := tokenManager.GetBaseTimeout() * time.Duration(1<<uint(i))
  ```
  This effectively doubles the backoff duration with each retry attempt.
- The backoff duration is capped at a maximum value:
  ```go
  if backoffDuration > tokenManager.GetMaxBackoffTime() {
      backoffDuration = tokenManager.GetMaxBackoffTime()
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
      float64(token.BaseTimeoutDuration)*math.Pow(2, float64(token.FailureCount-1)),
      float64(tm.maxBackoffTime),
  ))
  ```
- This calculation uses the `math.Pow` function to implement exponential growth based on the number of consecutive failures.
- The timeout duration is also capped at a maximum value (`tm.maxBackoffTime`).
- The token's `TimeoutUntil` is set to the current time plus this calculated duration.

This approach allows tokens that repeatedly fail increasingly longer "cool-down" periods before being used again, helping to manage rate limiting of individual tokens by the Pixiv API.

