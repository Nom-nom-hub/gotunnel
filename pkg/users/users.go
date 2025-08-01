package users

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
	"context"
	"sync"
)

// User represents a GoTunnel user
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	LastLogin   *time.Time `json:"last_login,omitempty"`
	APIKey      string    `json:"api_key,omitempty"`
	RateLimit   int       `json:"rate_limit"`
	MaxTunnels  int       `json:"max_tunnels"`
}

// UserManager handles user operations
type UserManager struct {
	db          *sql.DB
	jwtSecret   []byte
	rateLimiter *RateLimiter
}

// NewUserManager creates a new user manager
func NewUserManager(dbURL, jwtSecret string) (*UserManager, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &UserManager{
		db:          db,
		jwtSecret:   []byte(jwtSecret),
		rateLimiter: NewRateLimiter(),
	}, nil
}

// CreateUser creates a new user
func (um *UserManager) CreateUser(username, email, password string) (*User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate API key
	apiKey, err := um.generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Insert user
	var user User
	err = um.db.QueryRow(`
		INSERT INTO users (username, email, password_hash, api_key, role, status)
		VALUES ($1, $2, $3, $4, 'user', 'active')
		RETURNING id, username, email, role, status, created_at, rate_limit, max_tunnels
	`, username, email, string(hashedPassword), apiKey).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.Status,
		&user.CreatedAt, &user.RateLimit, &user.MaxTunnels,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.APIKey = apiKey
	return &user, nil
}

// AuthenticateUser authenticates a user
func (um *UserManager) AuthenticateUser(username, password string) (*User, error) {
	var user User
	var hashedPassword string

	err := um.db.QueryRow(`
		SELECT id, username, email, password_hash, role, status, created_at, rate_limit, max_tunnels
		FROM users WHERE username = $1 AND status = 'active'
	`, username).Scan(
		&user.ID, &user.Username, &user.Email, &hashedPassword, &user.Role, &user.Status,
		&user.CreatedAt, &user.RateLimit, &user.MaxTunnels,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	_, err = um.db.Exec(`UPDATE users SET last_login = NOW() WHERE id = $1`, user.ID)
	if err != nil {
		logrus.Errorf("Failed to update last login: %v", err)
	}

	return &user, nil
}

// GenerateToken generates a JWT token for a user
func (um *UserManager) GenerateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(um.jwtSecret)
}

// ValidateToken validates a JWT token
func (um *UserManager) ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return um.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return um.GetUserByID(userID)
	}

	return nil, fmt.Errorf("invalid token")
}

// GetUserByID gets a user by ID
func (um *UserManager) GetUserByID(id int) (*User, error) {
	var user User
	err := um.db.QueryRow(`
		SELECT id, username, email, role, status, created_at, last_login, rate_limit, max_tunnels
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.Status,
		&user.CreatedAt, &user.LastLogin, &user.RateLimit, &user.MaxTunnels,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// GetUserByAPIKey gets a user by API key
func (um *UserManager) GetUserByAPIKey(apiKey string) (*User, error) {
	var user User
	err := um.db.QueryRow(`
		SELECT id, username, email, role, status, created_at, last_login, rate_limit, max_tunnels
		FROM users WHERE api_key = $1 AND status = 'active'
	`, apiKey).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.Status,
		&user.CreatedAt, &user.LastLogin, &user.RateLimit, &user.MaxTunnels,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}

	return &user, nil
}

// UpdateUser updates a user
func (um *UserManager) UpdateUser(id int, updates map[string]interface{}) error {
	// Build dynamic query
	var setClauses []string
	var args []interface{}
	argIndex := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIndex)

	_, err := um.db.Exec(query, args...)
	return err
}

// DeleteUser deletes a user
func (um *UserManager) DeleteUser(id int) error {
	_, err := um.db.Exec("UPDATE users SET status = 'deleted' WHERE id = $1", id)
	return err
}

// ListUsers lists all users (admin only)
func (um *UserManager) ListUsers(limit, offset int) ([]*User, error) {
	rows, err := um.db.Query(`
		SELECT id, username, email, role, status, created_at, last_login, rate_limit, max_tunnels
		FROM users WHERE status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Role, &user.Status,
			&user.CreatedAt, &user.LastLogin, &user.RateLimit, &user.MaxTunnels,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// CheckRateLimit checks if a user has exceeded their rate limit
func (um *UserManager) CheckRateLimit(userID int, endpoint string) bool {
	return um.rateLimiter.CheckLimit(userID, endpoint)
}

// generateAPIKey generates a secure API key
func (um *UserManager) generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("gt_%x", bytes), nil
}

// RateLimiter handles rate limiting
type RateLimiter struct {
	limits map[string]map[int]time.Time
	mu     sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limits: make(map[string]map[int]time.Time),
	}
}

// CheckLimit checks if a user has exceeded their rate limit
func (rl *RateLimiter) CheckLimit(userID int, endpoint string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := fmt.Sprintf("%s:%d", endpoint, userID)
	now := time.Now()

	if userLimits, exists := rl.limits[key]; exists {
		if lastRequest, exists := userLimits[userID]; exists {
			if now.Sub(lastRequest) < time.Minute {
				return false // Rate limit exceeded
			}
		}
	}

	// Update last request time
	if rl.limits[key] == nil {
		rl.limits[key] = make(map[int]time.Time)
	}
	rl.limits[key][userID] = now
	return true
}

// Middleware functions for HTTP handlers

// AuthMiddleware validates JWT tokens
func (um *UserManager) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		user, err := um.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RateLimitMiddleware applies rate limiting
func (um *UserManager) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*User)
		endpoint := r.URL.Path

		if !um.CheckRateLimit(user.ID, endpoint) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// AdminMiddleware ensures user is admin
func (um *UserManager) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*User)
		if user.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
} 