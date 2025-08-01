package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

// Token represents an authentication token
type Token struct {
	ID        string
	Value     string
	CreatedAt time.Time
	ExpiresAt time.Time
	Active    bool
}

// TokenManager handles token generation and validation
type TokenManager struct {
	tokens map[string]*Token
}

// NewTokenManager creates a new token manager
func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokens: make(map[string]*Token),
	}
}

// GenerateToken creates a new authentication token
func (tm *TokenManager) GenerateToken(expiration time.Duration) (*Token, error) {
	// Generate random bytes for token
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Create token value (base64 encoded)
	tokenValue := base64.URLEncoding.EncodeToString(randomBytes)
	
	// Create token hash for storage
	hash := sha256.Sum256([]byte(tokenValue))
	tokenID := hex.EncodeToString(hash[:])

	now := time.Now()
	token := &Token{
		ID:        tokenID,
		Value:     tokenValue,
		CreatedAt: now,
		ExpiresAt: now.Add(expiration),
		Active:    true,
	}

	tm.tokens[tokenID] = token
	return token, nil
}

// ValidateToken validates a token and returns the token if valid
func (tm *TokenManager) ValidateToken(tokenValue string) (*Token, bool) {
	// Hash the token value
	hash := sha256.Sum256([]byte(tokenValue))
	tokenID := hex.EncodeToString(hash[:])

	token, exists := tm.tokens[tokenID]
	if !exists {
		return nil, false
	}

	// Check if token is active and not expired
	if !token.Active || time.Now().After(token.ExpiresAt) {
		return nil, false
	}

	return token, true
}

// RevokeToken revokes a token
func (tm *TokenManager) RevokeToken(tokenValue string) bool {
	hash := sha256.Sum256([]byte(tokenValue))
	tokenID := hex.EncodeToString(hash[:])

	token, exists := tm.tokens[tokenID]
	if !exists {
		return false
	}

	token.Active = false
	return true
}

// ListTokens returns all active tokens
func (tm *TokenManager) ListTokens() []*Token {
	var activeTokens []*Token
	now := time.Now()

	for _, token := range tm.tokens {
		if token.Active && now.Before(token.ExpiresAt) {
			activeTokens = append(activeTokens, token)
		}
	}

	return activeTokens
}

// CleanupExpired removes expired tokens
func (tm *TokenManager) CleanupExpired() int {
	now := time.Now()
	removed := 0

	for id, token := range tm.tokens {
		if now.After(token.ExpiresAt) {
			delete(tm.tokens, id)
			removed++
		}
	}

	return removed
}

// SimpleAuth provides a simple authentication mechanism
type SimpleAuth struct {
	tokenManager *TokenManager
	allowedTokens map[string]bool
}

// NewSimpleAuth creates a new simple authentication handler
func NewSimpleAuth() *SimpleAuth {
	return &SimpleAuth{
		tokenManager:  NewTokenManager(),
		allowedTokens: make(map[string]bool),
	}
}

// AddAllowedToken adds a token to the allowed list
func (sa *SimpleAuth) AddAllowedToken(token string) {
	sa.allowedTokens[token] = true
}

// Authenticate validates a token
func (sa *SimpleAuth) Authenticate(token string) bool {
	// Check if token is in allowed list
	if sa.allowedTokens[token] {
		return true
	}

	// Check if token is valid in token manager
	_, valid := sa.tokenManager.ValidateToken(token)
	return valid
}

// GenerateClientToken generates a token for client use
func (sa *SimpleAuth) GenerateClientToken(expiration time.Duration) (string, error) {
	token, err := sa.tokenManager.GenerateToken(expiration)
	if err != nil {
		return "", err
	}
	return token.Value, nil
} 