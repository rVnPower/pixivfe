package token_manager

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type TokenStatus int

const (
	Good TokenStatus = iota
	TimedOut
)

type Token struct {
	Value               string
	Status              TokenStatus
	TimeoutUntil        time.Time
	FailureCount        int
	LastUsed            time.Time
	BaseTimeoutDuration time.Duration
}

type TokenManager struct {
	tokens              []*Token
	mu                  sync.Mutex
	maxRetries          int
	baseTimeout         time.Duration
	maxBackoffTime      time.Duration
	loadBalancingMethod string
	currentIndex        int
}

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

func (tm *TokenManager) getHealthyTokens() []*Token {
	healthyTokens := make([]*Token, 0)
	for _, token := range tm.tokens {
		if token.Status == Good {
			healthyTokens = append(healthyTokens, token)
		}
	}
	return healthyTokens
}

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

func (tm *TokenManager) roundRobinSelection(healthyTokens []*Token) *Token {
	if tm.currentIndex >= len(healthyTokens) {
		tm.currentIndex = 0
	}
	selectedToken := healthyTokens[tm.currentIndex]
	tm.currentIndex++
	return selectedToken
}

func (tm *TokenManager) randomSelection(healthyTokens []*Token) *Token {
	return healthyTokens[rand.Intn(len(healthyTokens))]
}

func (tm *TokenManager) leastRecentlyUsedSelection(healthyTokens []*Token) *Token {
	sort.Slice(healthyTokens, func(i, j int) bool {
		return healthyTokens[i].LastUsed.Before(healthyTokens[j].LastUsed)
	})
	return healthyTokens[0]
}

func (tm *TokenManager) MarkTokenStatus(token *Token, status TokenStatus) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	token.Status = status
	if status == TimedOut {
		token.FailureCount++
		timeoutDuration := time.Duration(math.Min(
			float64(token.BaseTimeoutDuration)*math.Pow(2, float64(token.FailureCount-1)),
			float64(tm.maxBackoffTime),
		))
		token.TimeoutUntil = time.Now().Add(timeoutDuration)
	} else {
		token.FailureCount = 0
		token.BaseTimeoutDuration = time.Second
	}
}

func (tm *TokenManager) ResetAllTokens() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, token := range tm.tokens {
		token.Status = Good
		token.FailureCount = 0
		token.BaseTimeoutDuration = time.Second
	}
}

func (tm *TokenManager) GetMaxRetries() int {
	return tm.maxRetries
}

func (tm *TokenManager) GetBaseTimeout() time.Duration {
	return tm.baseTimeout
}

func (tm *TokenManager) GetMaxBackoffTime() time.Duration {
	return tm.maxBackoffTime
}
