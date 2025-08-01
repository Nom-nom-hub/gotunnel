package billing

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
)

// Plan represents a subscription plan
type Plan struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	Currency    string                 `json:"currency"`
	BillingCycle string                `json:"billing_cycle"` // monthly, yearly
	Features    map[string]interface{} `json:"features"`
	Limits      map[string]int         `json:"limits"`
	IsActive    bool                   `json:"is_active"`
}

// Subscription represents a user's subscription
type Subscription struct {
	ID            string    `json:"id"`
	UserID        int       `json:"user_id"`
	PlanID        string    `json:"plan_id"`
	Status        string    `json:"status"` // active, canceled, expired
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	AutoRenew     bool      `json:"auto_renew"`
	PaymentMethod string    `json:"payment_method"`
	LastBilling   time.Time `json:"last_billing"`
	NextBilling   time.Time `json:"next_billing"`
}

// BillingManager handles billing and subscription operations
type BillingManager struct {
	db *sql.DB
}

// NewBillingManager creates a new billing manager
func NewBillingManager(dbURL string) (*BillingManager, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &BillingManager{db: db}, nil
}

// GetPlans returns all available plans
func (bm *BillingManager) GetPlans() ([]*Plan, error) {
	rows, err := bm.db.Query(`
		SELECT id, name, description, price, currency, billing_cycle, features, limits, is_active
		FROM plans WHERE is_active = true ORDER BY price ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*Plan
	for rows.Next() {
		var plan Plan
		var featuresJSON, limitsJSON []byte

		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description, &plan.Price, &plan.Currency,
			&plan.BillingCycle, &featuresJSON, &limitsJSON, &plan.IsActive,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON features and limits
		if err := json.Unmarshal(featuresJSON, &plan.Features); err != nil {
			logrus.Errorf("Failed to parse plan features: %v", err)
		}
		if err := json.Unmarshal(limitsJSON, &plan.Limits); err != nil {
			logrus.Errorf("Failed to parse plan limits: %v", err)
		}

		plans = append(plans, &plan)
	}

	return plans, nil
}

// GetUserSubscription gets a user's current subscription
func (bm *BillingManager) GetUserSubscription(userID int) (*Subscription, error) {
	var sub Subscription
	err := bm.db.QueryRow(`
		SELECT id, user_id, plan_id, status, start_date, end_date, auto_renew, 
		       payment_method, last_billing, next_billing
		FROM subscriptions 
		WHERE user_id = $1 AND status = 'active'
		ORDER BY start_date DESC LIMIT 1
	`, userID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.Status, &sub.StartDate, &sub.EndDate,
		&sub.AutoRenew, &sub.PaymentMethod, &sub.LastBilling, &sub.NextBilling,
	)

	if err != nil {
		return nil, fmt.Errorf("no active subscription found")
	}

	return &sub, nil
}

// CreateSubscription creates a new subscription
func (bm *BillingManager) CreateSubscription(userID int, planID string, paymentMethod string) (*Subscription, error) {
	// Get plan details
	var plan Plan
	var featuresJSON, limitsJSON []byte
	err := bm.db.QueryRow(`
		SELECT id, name, description, price, currency, billing_cycle, features, limits
		FROM plans WHERE id = $1 AND is_active = true
	`, planID).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.Price, &plan.Currency,
		&plan.BillingCycle, &featuresJSON, &limitsJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("plan not found")
	}

	// Calculate subscription dates
	now := time.Now()
	var endDate time.Time
	if plan.BillingCycle == "yearly" {
		endDate = now.AddDate(1, 0, 0)
	} else {
		endDate = now.AddDate(0, 1, 0)
	}

	// Create subscription
	var subID string
	err = bm.db.QueryRow(`
		INSERT INTO subscriptions (user_id, plan_id, status, start_date, end_date, 
		                         auto_renew, payment_method, last_billing, next_billing)
		VALUES ($1, $2, 'active', $3, $4, true, $5, $3, $4)
		RETURNING id
	`, userID, planID, now, endDate, paymentMethod).Scan(&subID)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Update user's plan limits
	if err := bm.updateUserLimits(userID, plan.Limits); err != nil {
		logrus.Errorf("Failed to update user limits: %v", err)
	}

	return &Subscription{
		ID:            subID,
		UserID:        userID,
		PlanID:        planID,
		Status:        "active",
		StartDate:     now,
		EndDate:       endDate,
		AutoRenew:     true,
		PaymentMethod: paymentMethod,
		LastBilling:   now,
		NextBilling:   endDate,
	}, nil
}

// CancelSubscription cancels a user's subscription
func (bm *BillingManager) CancelSubscription(userID int) error {
	_, err := bm.db.Exec(`
		UPDATE subscriptions 
		SET status = 'canceled', auto_renew = false 
		WHERE user_id = $1 AND status = 'active'
	`, userID)
	return err
}

// ProcessBilling processes billing for all active subscriptions
func (bm *BillingManager) ProcessBilling() error {
	rows, err := bm.db.Query(`
		SELECT s.id, s.user_id, s.plan_id, s.next_billing, p.price, p.billing_cycle
		FROM subscriptions s
		JOIN plans p ON s.plan_id = p.id
		WHERE s.status = 'active' AND s.next_billing <= NOW()
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var subID string
		var userID int
		var planID string
		var nextBilling time.Time
		var price float64
		var billingCycle string

		if err := rows.Scan(&subID, &userID, &planID, &nextBilling, &price, &billingCycle); err != nil {
			continue
		}

		// Process billing for this subscription
		if err := bm.processSubscriptionBilling(subID, userID, planID, price, billingCycle); err != nil {
			logrus.Errorf("Failed to process billing for subscription %s: %v", subID, err)
		}
	}

	return nil
}

// processSubscriptionBilling processes billing for a single subscription
func (bm *BillingManager) processSubscriptionBilling(subID string, userID int, planID string, price float64, billingCycle string) error {
	// Calculate next billing date
	var nextBilling time.Time
	if billingCycle == "yearly" {
		nextBilling = time.Now().AddDate(1, 0, 0)
	} else {
		nextBilling = time.Now().AddDate(0, 1, 0)
	}

	// Update subscription billing dates
	_, err := bm.db.Exec(`
		UPDATE subscriptions 
		SET last_billing = NOW(), next_billing = $1
		WHERE id = $2
	`, nextBilling, subID)

	return err
}

// updateUserLimits updates user limits based on plan
func (bm *BillingManager) updateUserLimits(userID int, limits map[string]int) error {
	// Update user rate limits and tunnel limits
	if maxTunnels, ok := limits["max_tunnels"]; ok {
		_, err := bm.db.Exec(`UPDATE users SET max_tunnels = $1 WHERE id = $2`, maxTunnels, userID)
		if err != nil {
			return err
		}
	}

	if rateLimit, ok := limits["rate_limit"]; ok {
		_, err := bm.db.Exec(`UPDATE users SET rate_limit = $1 WHERE id = $2`, rateLimit, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckUsage checks if user has exceeded their plan limits
func (bm *BillingManager) CheckUsage(userID int, resource string) (bool, error) {
	// Get user's current subscription
	sub, err := bm.GetUserSubscription(userID)
	if err != nil {
		// No subscription, use free tier limits
		return bm.checkFreeTierUsage(userID, resource)
	}

	// Get plan limits
	var limitsJSON []byte
	err = bm.db.QueryRow(`SELECT limits FROM plans WHERE id = $1`, sub.PlanID).Scan(&limitsJSON)
	if err != nil {
		return false, err
	}

	var limits map[string]int
	if err := json.Unmarshal(limitsJSON, &limits); err != nil {
		return false, err
	}

	// Check specific resource usage
	switch resource {
	case "tunnels":
		return bm.checkTunnelUsage(userID, limits["max_tunnels"])
	case "bandwidth":
		return bm.checkBandwidthUsage(userID, limits["bandwidth_limit"])
	case "requests":
		return bm.checkRequestUsage(userID, limits["rate_limit"])
	default:
		return true, nil
	}
}

// checkFreeTierUsage checks usage for free tier users
func (bm *BillingManager) checkFreeTierUsage(userID int, resource string) (bool, error) {
	freeLimits := map[string]int{
		"max_tunnels":    3,
		"bandwidth_limit": 100 * 1024 * 1024, // 100MB
		"rate_limit":      100,
	}

	switch resource {
	case "tunnels":
		return bm.checkTunnelUsage(userID, freeLimits["max_tunnels"])
	case "bandwidth":
		return bm.checkBandwidthUsage(userID, freeLimits["bandwidth_limit"])
	case "requests":
		return bm.checkRequestUsage(userID, freeLimits["rate_limit"])
	default:
		return true, nil
	}
}

// checkTunnelUsage checks if user has exceeded tunnel limit
func (bm *BillingManager) checkTunnelUsage(userID, maxTunnels int) (bool, error) {
	var count int
	err := bm.db.QueryRow(`
		SELECT COUNT(*) FROM tunnels WHERE user_id = $1 AND status = 'active'
	`, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count < maxTunnels, nil
}

// checkBandwidthUsage checks if user has exceeded bandwidth limit
func (bm *BillingManager) checkBandwidthUsage(userID int, bandwidthLimit int) (bool, error) {
	var totalBytes int64
	err := bm.db.QueryRow(`
		SELECT COALESCE(SUM(bytes_transferred), 0) FROM tunnels WHERE user_id = $1
	`, userID).Scan(&totalBytes)
	if err != nil {
		return false, err
	}

	return int(totalBytes) < bandwidthLimit, nil
}

// checkRequestUsage checks if user has exceeded request rate limit
func (bm *BillingManager) checkRequestUsage(userID int, rateLimit int) (bool, error) {
	// This would typically check against a rate limiting system
	// For now, we'll return true (allowed)
	return true, nil
}

// GetUsageStats gets usage statistics for a user
func (bm *BillingManager) GetUsageStats(userID int) (map[string]interface{}, error) {
	// Get current subscription
	sub, err := bm.GetUserSubscription(userID)
	if err != nil {
		// Return free tier stats
		return bm.getFreeTierStats(userID)
	}

	// Get plan details
	var plan Plan
	var featuresJSON, limitsJSON []byte
	err = bm.db.QueryRow(`
		SELECT id, name, description, price, currency, billing_cycle, features, limits
		FROM plans WHERE id = $1
	`, sub.PlanID).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.Price, &plan.Currency,
		&plan.BillingCycle, &featuresJSON, &limitsJSON,
	)
	if err != nil {
		return nil, err
	}

	// Parse limits
	if err := json.Unmarshal(limitsJSON, &plan.Limits); err != nil {
		return nil, err
	}

	// Get current usage
	stats, err := bm.getCurrentUsage(userID)
	if err != nil {
		return nil, err
	}

	// Add plan information
	stats["plan"] = plan
	stats["subscription"] = sub

	return stats, nil
}

// getFreeTierStats gets usage stats for free tier users
func (bm *BillingManager) getFreeTierStats(userID int) (map[string]interface{}, error) {
	stats, err := bm.getCurrentUsage(userID)
	if err != nil {
		return nil, err
	}

	// Add free tier limits
	stats["limits"] = map[string]int{
		"max_tunnels":    3,
		"bandwidth_limit": 100 * 1024 * 1024, // 100MB
		"rate_limit":      100,
	}

	stats["plan"] = map[string]interface{}{
		"name":        "Free",
		"price":       0.0,
		"billing_cycle": "none",
	}

	return stats, nil
}

// getCurrentUsage gets current usage statistics
func (bm *BillingManager) getCurrentUsage(userID int) (map[string]interface{}, error) {
	var activeTunnels, totalTunnels int
	var totalBytes int64

	// Get tunnel counts
	err := bm.db.QueryRow(`
		SELECT 
			COUNT(*) as total_tunnels,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_tunnels,
			COALESCE(SUM(bytes_transferred), 0) as total_bytes
		FROM tunnels WHERE user_id = $1
	`, userID).Scan(&totalTunnels, &activeTunnels, &totalBytes)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"active_tunnels": activeTunnels,
		"total_tunnels":  totalTunnels,
		"total_bytes":    totalBytes,
		"total_mb":       totalBytes / 1024 / 1024,
	}, nil
} 