// Package token_manager provides functionality for managing and rotating API tokens
// with features like load balancing, timeout handling, and exponential backoff.
package token_manager

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// TokenStatus represents the current state of a token
type TokenStatus int

const (
	Good     TokenStatus = iota // Token is in a good state and can be used
	TimedOut                    // Token is currently timed out and should not be used
)

// Token represents an individual API token with its associated metadata
type Token struct {
	Value               string        // The actual token value
	Status              TokenStatus   // Current status of the token
	TimeoutUntil        time.Time     // Time until which the token is timed out
	FailureCount        int           // Number of consecutive failures
	LastUsed            time.Time     // Last time the token was used
	BaseTimeoutDuration time.Duration // Base duration for timeout calculations
}

// TokenManager handles a collection of tokens and provides methods for token selection and management
type TokenManager struct {
	tokens              []*Token      // Slice of available tokens
	mu                  sync.Mutex    // Mutex for thread-safe operations
	maxRetries          int           // Maximum number of retries before considering a request failed
	baseTimeout         time.Duration // Base timeout duration for requests
	maxBackoffTime      time.Duration // Maximum allowed backoff time
	loadBalancingMethod string        // Method used for load balancing (e.g., "round-robin", "random")
	currentIndex        int           // Current index for round-robin selection
}

// NewTokenManager creates and initializes a new TokenManager with the given parameters
func NewTokenManager(tokenValues []string, maxRetries int, baseTimeout, maxBackoffTime time.Duration, loadBalancingMethod string) *TokenManager {
	tokens := make([]*Token, len(tokenValues))
	for i, value := range tokenValues {
		tokens[i] = &Token{
			Value:               value,
			Status:              Good,
			BaseTimeoutDuration: time.Second,
		}
	}
	return &TokenManager{
		tokens:              tokens,
		maxRetries:          maxRetries,
		baseTimeout:         baseTimeout,
		maxBackoffTime:      maxBackoffTime,
		loadBalancingMethod: loadBalancingMethod,
		currentIndex:        0,
	}
}

// GetToken selects and returns a token based on the configured load balancing method
func (tm *TokenManager) GetToken() *Token {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	now := time.Now()
	healthyTokens := tm.getHealthyTokens()

	if len(healthyTokens) == 0 {
		return tm.getFallbackToken(now)
	}

	var selectedToken *Token
	switch tm.loadBalancingMethod {
	case "round-robin":
		selectedToken = tm.roundRobinSelection(healthyTokens)
	case "random":
		selectedToken = tm.randomSelection(healthyTokens)
	case "least-recently-used":
		selectedToken = tm.leastRecentlyUsedSelection(healthyTokens)
	default:
		selectedToken = tm.roundRobinSelection(healthyTokens)
	}

	selectedToken.LastUsed = now
	return selectedToken
}

// getHealthyTokens returns a slice of tokens that are currently in a good state
func (tm *TokenManager) getHealthyTokens() []*Token {
	healthyTokens := make([]*Token, 0)
	for _, token := range tm.tokens {
		if token.Status == Good {
			healthyTokens = append(healthyTokens, token)
		}
	}
	return healthyTokens
}

// getFallbackToken attempts to find a timed-out token that can be reset and used
func (tm *TokenManager) getFallbackToken(now time.Time) *Token {
	var bestToken *Token
	for _, token := range tm.tokens {
		if token.Status == TimedOut && (bestToken == nil || token.TimeoutUntil.Before(bestToken.TimeoutUntil)) {
			bestToken = token
		}
	}

	if bestToken != nil && now.After(bestToken.TimeoutUntil) {
		bestToken.Status = Good
		bestToken.LastUsed = now
		return bestToken
	}

	return bestToken
}

// roundRobinSelection implements the round-robin token selection strategy
func (tm *TokenManager) roundRobinSelection(healthyTokens []*Token) *Token {
	if tm.currentIndex >= len(healthyTokens) {
		tm.currentIndex = 0
	}
	selectedToken := healthyTokens[tm.currentIndex]
	tm.currentIndex++
	return selectedToken
}

// randomSelection implements the random token selection strategy
func (tm *TokenManager) randomSelection(healthyTokens []*Token) *Token {
	return healthyTokens[rand.Intn(len(healthyTokens))]
}

// leastRecentlyUsedSelection implements the least recently used token selection strategy
func (tm *TokenManager) leastRecentlyUsedSelection(healthyTokens []*Token) *Token {
	sort.Slice(healthyTokens, func(i, j int) bool {
		return healthyTokens[i].LastUsed.Before(healthyTokens[j].LastUsed)
	})
	return healthyTokens[0]
}

// MarkTokenStatus updates the status of a token and handles timeout logic
func (tm *TokenManager) MarkTokenStatus(token *Token, status TokenStatus) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	token.Status = status
	if status == TimedOut {
		token.FailureCount++
		// Calculate timeout duration using exponential backoff with a maximum limit
		timeoutDuration := time.Duration(math.Min(
			float64(token.BaseTimeoutDuration)*math.Pow(2, float64(token.FailureCount-1)),
			float64(tm.maxBackoffTime),
		))
		token.TimeoutUntil = time.Now().Add(timeoutDuration)
	} else {
		// Reset failure count and base timeout duration when marked as Good
		token.FailureCount = 0
		token.BaseTimeoutDuration = time.Second
	}
}

// ResetAllTokens resets all tokens to their initial good state
func (tm *TokenManager) ResetAllTokens() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, token := range tm.tokens {
		token.Status = Good
		token.FailureCount = 0
		token.BaseTimeoutDuration = time.Second
	}
}

// GetMaxRetries returns the maximum number of retries configured for the TokenManager
func (tm *TokenManager) GetMaxRetries() int {
	return tm.maxRetries
}

// GetBaseTimeout returns the base timeout duration configured for the TokenManager
func (tm *TokenManager) GetBaseTimeout() time.Duration {
	return tm.baseTimeout
}

// GetMaxBackoffTime returns the maximum backoff time configured for the TokenManager
func (tm *TokenManager) GetMaxBackoffTime() time.Duration {
	return tm.maxBackoffTime
}
