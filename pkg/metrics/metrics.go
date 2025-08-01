package metrics

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
)

// Metrics represents the metrics collection system
type Metrics struct {
	mu sync.RWMutex

	// Prometheus metrics
	activeTunnels     prometheus.Gauge
	totalConnections  prometheus.Counter
	bytesTransferred  prometheus.Counter
	requestDuration   prometheus.Histogram
	errorsTotal       prometheus.Counter
	requestsTotal     prometheus.Counter

	// Database connection
	db *sql.DB

	// Configuration
	enabled bool
	port    string
}

// NewMetrics creates a new metrics collector
func NewMetrics(dbURL, port string) (*Metrics, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	m := &Metrics{
		db:      db,
		enabled: true,
		port:    port,
	}

	// Initialize Prometheus metrics
	m.activeTunnels = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gotunnel_active_tunnels",
		Help: "Number of active tunnels",
	})

	m.totalConnections = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gotunnel_total_connections",
		Help: "Total number of connections",
	})

	m.bytesTransferred = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gotunnel_bytes_transferred",
		Help: "Total bytes transferred",
	})

	m.requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "gotunnel_request_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.DefBuckets,
	})

	m.errorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gotunnel_errors_total",
		Help: "Total number of errors",
	})

	m.requestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gotunnel_requests_total",
		Help: "Total number of requests",
	})

	// Register metrics
	prometheus.MustRegister(
		m.activeTunnels,
		m.totalConnections,
		m.bytesTransferred,
		m.requestDuration,
		m.errorsTotal,
		m.requestsTotal,
	)

	return m, nil
}

// Start starts the metrics server
func (m *Metrics) Start(ctx context.Context) error {
	if !m.enabled {
		return nil
	}

	// Start metrics HTTP server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", m.healthHandler)
	mux.HandleFunc("/stats", m.statsHandler)

	server := &http.Server{
		Addr:    ":" + m.port,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("Metrics server error: %v", err)
		}
	}()

	// Start background metrics collection
	go m.collectMetrics(ctx)

	logrus.Infof("Metrics server started on port %s", m.port)
	return nil
}

// RecordTunnelCreated records a new tunnel creation
func (m *Metrics) RecordTunnelCreated(subdomain string, userID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.activeTunnels.Inc()

	// Store in database
	_, err := m.db.Exec(`
		INSERT INTO metrics (tunnel_id, metric_type, metric_value, metric_unit, metadata)
		SELECT id, 'tunnel_created', 1, 'count', $1::jsonb
		FROM tunnels WHERE subdomain = $2
	`, map[string]interface{}{
		"subdomain": subdomain,
		"user_id":   userID,
		"timestamp": time.Now(),
	}, subdomain)

	if err != nil {
		logrus.Errorf("Failed to record tunnel creation: %v", err)
	}
}

// RecordTunnelClosed records a tunnel closure
func (m *Metrics) RecordTunnelClosed(subdomain string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.activeTunnels.Dec()

	// Update database
	_, err := m.db.Exec(`
		UPDATE tunnels SET status = 'inactive', updated_at = NOW()
		WHERE subdomain = $1
	`, subdomain)

	if err != nil {
		logrus.Errorf("Failed to record tunnel closure: %v", err)
	}
}

// RecordConnection records a new connection
func (m *Metrics) RecordConnection(subdomain string, clientIP string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalConnections.Inc()

	// Store session
	_, err := m.db.Exec(`
		INSERT INTO tunnel_sessions (tunnel_id, session_id, client_ip, started_at)
		SELECT id, $1, $2, NOW()
		FROM tunnels WHERE subdomain = $3
	`, fmt.Sprintf("%s-%d", subdomain, time.Now().Unix()), clientIP, subdomain)

	if err != nil {
		logrus.Errorf("Failed to record connection: %v", err)
	}
}

// RecordBytesTransferred records bytes transferred
func (m *Metrics) RecordBytesTransferred(subdomain string, bytes int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.bytesTransferred.Add(float64(bytes))

	// Update database
	_, err := m.db.Exec(`
		UPDATE tunnels 
		SET bytes_transferred = bytes_transferred + $1, 
		    last_activity = NOW(),
		    updated_at = NOW()
		WHERE subdomain = $2
	`, bytes, subdomain)

	if err != nil {
		logrus.Errorf("Failed to record bytes transferred: %v", err)
	}
}

// RecordRequest records a request with duration
func (m *Metrics) RecordRequest(subdomain string, duration time.Duration, statusCode int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestsTotal.Inc()
	m.requestDuration.Observe(duration.Seconds())

	if statusCode >= 400 {
		m.errorsTotal.Inc()
	}

	// Store metric
	_, err := m.db.Exec(`
		INSERT INTO metrics (tunnel_id, metric_type, metric_value, metric_unit, metadata)
		SELECT id, 'request', 1, 'count', $1::jsonb
		FROM tunnels WHERE subdomain = $2
	`, map[string]interface{}{
		"duration":   duration.Milliseconds(),
		"status_code": statusCode,
		"timestamp":   time.Now(),
	}, subdomain)

	if err != nil {
		logrus.Errorf("Failed to record request: %v", err)
	}
}

// healthHandler handles health check requests
func (m *Metrics) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"time":   time.Now(),
	})
}

// statsHandler provides detailed statistics
func (m *Metrics) statsHandler(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Get database stats
	var stats struct {
		ActiveTunnels    int64   `json:"active_tunnels"`
		TotalConnections int64   `json:"total_connections"`
		BytesTransferred int64   `json:"bytes_transferred"`
		AvgResponseTime  float64 `json:"avg_response_time"`
		ErrorRate        float64 `json:"error_rate"`
	}

	// Query database for stats
	err := m.db.QueryRow(`
		SELECT 
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_tunnels,
			SUM(connection_count) as total_connections,
			SUM(bytes_transferred) as bytes_transferred
		FROM tunnels
	`).Scan(&stats.ActiveTunnels, &stats.TotalConnections, &stats.BytesTransferred)

	if err != nil {
		logrus.Errorf("Failed to get stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// collectMetrics periodically collects and stores metrics
func (m *Metrics) collectMetrics(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.collectSystemMetrics()
		}
	}
}

// collectSystemMetrics collects system-level metrics
func (m *Metrics) collectSystemMetrics() {
	// Collect user activity metrics
	rows, err := m.db.Query(`
		SELECT 
			u.username,
			COUNT(t.id) as tunnel_count,
			SUM(t.bytes_transferred) as total_bytes
		FROM users u
		LEFT JOIN tunnels t ON u.id = t.user_id
		GROUP BY u.id, u.username
		ORDER BY total_bytes DESC
		LIMIT 10
	`)

	if err != nil {
		logrus.Errorf("Failed to collect user metrics: %v", err)
		return
	}
	defer rows.Close()

	// Process and store user metrics
	for rows.Next() {
		var username string
		var tunnelCount int
		var totalBytes int64

		if err := rows.Scan(&username, &tunnelCount, &totalBytes); err != nil {
			logrus.Errorf("Failed to scan user metrics: %v", err)
			continue
		}

		// Store user activity metric
		_, err := m.db.Exec(`
			INSERT INTO metrics (tunnel_id, metric_type, metric_value, metric_unit, metadata)
			SELECT t.id, 'user_activity', $1, 'bytes', $2::jsonb
			FROM tunnels t
			JOIN users u ON t.user_id = u.id
			WHERE u.username = $3
			LIMIT 1
		`, float64(totalBytes), map[string]interface{}{
			"username":     username,
			"tunnel_count": tunnelCount,
			"timestamp":    time.Now(),
		}, username)

		if err != nil {
			logrus.Errorf("Failed to store user activity: %v", err)
		}
	}
}

// Close closes the metrics system
func (m *Metrics) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
} 