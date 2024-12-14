// pkg/metrics/collector.go
package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"log"
	"net/http"
	"sync"
)

// RedditMetrics holds all metrics for the Reddit clone
type RedditMetrics struct {
	PostsCreated        prometheus.Counter
	CommentsCreated     prometheus.Counter
	VotesRecorded       prometheus.Counter
	ActiveUsers         prometheus.Gauge
	ResponseTime        prometheus.Histogram
	SubredditSize       *prometheus.GaugeVec
	TotalUsers          prometheus.Counter
	ErrorCount          prometheus.Counter
	PersonaStats        map[string]*PersonaStats
	SimulatedUsers      prometheus.Gauge
	AverageResponseTime prometheus.Gauge
	ErrorRate           prometheus.Gauge
}

type PersonaStats struct {
	ActiveUsers  int
	PostCount    int
	VoteCount    int
	CommentCount int
}

var (
	metricsOnce     sync.Once
	metricsInstance *RedditMetrics
)

// NewRedditMetrics creates a new Reddit metrics collector
func NewRedditMetrics() *RedditMetrics {
	metricsOnce.Do(func() {
		metricsInstance = &RedditMetrics{
			PostsCreated: promauto.NewCounter(prometheus.CounterOpts{
				Name: "reddit_posts_total",
				Help: "Total number of posts created",
			}),
			CommentsCreated: promauto.NewCounter(prometheus.CounterOpts{
				Name: "reddit_comments_total",
				Help: "Total number of comments created",
			}),
			VotesRecorded: promauto.NewCounter(prometheus.CounterOpts{
				Name: "reddit_votes_total",
				Help: "Total number of votes recorded",
			}),
			ActiveUsers: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "reddit_active_users",
				Help: "Number of currently active users",
			}),
			ResponseTime: promauto.NewHistogram(prometheus.HistogramOpts{
				Name:    "reddit_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			}),
			SubredditSize: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: "reddit_subreddit_members",
				Help: "Number of members per subreddit",
			}, []string{"subreddit"}),
			TotalUsers: promauto.NewCounter(prometheus.CounterOpts{
				Name: "reddit_total_users",
				Help: "Total number of registered users",
			}),
			ErrorCount: promauto.NewCounter(prometheus.CounterOpts{
				Name: "reddit_errors_total",
				Help: "Total number of errors",
			}),
			PersonaStats: make(map[string]*PersonaStats),
			SimulatedUsers: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "reddit_simulated_users_total",
				Help: "Total number of simulated users",
			}),
			AverageResponseTime: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "reddit_average_response_time_seconds",
				Help: "Average response time in seconds",
			}),
			ErrorRate: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "reddit_error_rate",
				Help: "Rate of errors per second",
			}),
		}

	})
	return metricsInstance
}

//delete this later
//func (m *RedditMetrics) Get() map[string]float64 {
//	metrics := make(map[string]float64)
//
//	// Get values using prometheus.Gauge methods
//	metrics["active_users"] = getGaugeValue(m.ActiveUsers)
//	metrics["simulated_users"] = getGaugeValue(m.SimulatedUsers)
//
//	// Get values using prometheus.Counter methods
//	metrics["total_users"] = getCounterValue(m.TotalUsers)
//	metrics["posts"] = getCounterValue(m.PostsCreated)
//	metrics["comments"] = getCounterValue(m.CommentsCreated)
//	metrics["votes"] = getCounterValue(m.VotesRecorded)
//
//	return metrics
//}
//
//// Helper functions to get values from different metric types
//func getGaugeValue(g prometheus.Gauge) float64 {
//	var m dto.Metric
//	g.Write(&m)
//	return m.GetGauge().GetValue()
//}
//
//func getCounterValue(c prometheus.Counter) float64 {
//	var m dto.Metric
//	c.Write(&m)
//	return m.GetCounter().GetValue()
//}

// RecordAction records metrics for different types of actions
func (m *RedditMetrics) RecordAction(persona string, actionType string) {
	if _, exists := m.PersonaStats[persona]; !exists {
		m.PersonaStats[persona] = &PersonaStats{}
	}

	switch actionType {
	case "post":
		m.PostsCreated.Inc()
		m.PersonaStats[persona].PostCount++
	case "comment":
		m.CommentsCreated.Inc()
		m.PersonaStats[persona].CommentCount++
	case "vote":
		m.VotesRecorded.Inc()
		m.PersonaStats[persona].VoteCount++
	}
}

// UpdatePersonaCount updates the active users count for a persona
func (m *RedditMetrics) UpdatePersonaCount(persona string, count int) {
	if _, exists := m.PersonaStats[persona]; !exists {
		m.PersonaStats[persona] = &PersonaStats{}
	}
	m.PersonaStats[persona].ActiveUsers = count
}

// UpdateSimulationMetrics updates general simulation metrics
func (m *RedditMetrics) UpdateSimulationMetrics(users, connectionRate float64) {
	m.ActiveUsers.Set(users)
}

// RecordSimulatedAction records the duration of a simulated action
func (m *RedditMetrics) RecordSimulatedAction(duration float64) {
	m.ResponseTime.Observe(duration)
	m.AverageResponseTime.Set(duration)
}

// UpdateSubredditMembers updates the member count for a subreddit
func (m *RedditMetrics) UpdateSubredditMembers(subredditID string, count float64) {
	m.SubredditSize.WithLabelValues(subredditID).Set(count)
}

// RecordError increments the error counter
func (m *RedditMetrics) RecordError() {
	m.ErrorCount.Inc()
	// Log the error for debugging
	log.Println("An error occurred")
}

// RecordRequest records the duration of a request
func (m *RedditMetrics) RecordRequest(duration float64) {
	m.ResponseTime.Observe(duration)
}

// ReportCurrentStats generates a report of current statistics
func (m *RedditMetrics) ReportCurrentStats() {
	// Implementation for reporting current stats
	// This can be customized based on your needs
}
func (m *RedditMetrics) RecordLoadTestMetrics(users, responseTime, errorRate float64) {
	m.SimulatedUsers.Set(users)
	m.AverageResponseTime.Set(responseTime)
	m.ErrorRate.Set(errorRate)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
func (m *RedditMetrics) UpdateActiveUsers(delta float64) {
	var metric dto.Metric

	// Get current ActiveUsers value
	err := m.ActiveUsers.Write(&metric)
	if err != nil {
		fmt.Printf("Error reading ActiveUsers metric: %v\n", err)
		return
	}
	currentActive := metric.GetGauge().GetValue()

	// Calculate new ActiveUsers value
	newActive := currentActive + delta
	if newActive < 0 {
		newActive = 0
	}
	m.ActiveUsers.Set(newActive)

	// Get current SimulatedUsers value
	err = m.SimulatedUsers.Write(&metric)
	if err != nil {
		fmt.Printf("Error reading SimulatedUsers metric: %v\n", err)
		return
	}
	currentSimulated := metric.GetGauge().GetValue()

	// Calculate new SimulatedUsers value
	newSimulated := currentSimulated + delta
	if newSimulated < 0 {
		newSimulated = 0
	}
	m.SimulatedUsers.Set(newSimulated)

	// Increment TotalUsers only if delta > 0 (since it's a counter)
	if delta > 0 {
		m.TotalUsers.Add(delta)
	}
}

func (m *RedditMetrics) GetEssentialMetrics() (map[string]float64, error) {
	metrics := make(map[string]float64)
	var dtoMetric dto.Metric // Use a correctly typed dto.Metric object

	// Helper function to write metric and handle error
	writeMetric := func(promMetric prometheus.Metric, key string, getValue func(*dto.Metric) float64) error {
		// Write the metric value into dtoMetric
		err := promMetric.Write(&dtoMetric)
		if err != nil {
			return fmt.Errorf("error writing %s metric: %w", key, err)
		}
		metrics[key] = getValue(&dtoMetric)
		return nil
	}

	// Get Gauge values
	if err := writeMetric(m.ActiveUsers, "active_users", func(m *dto.Metric) float64 { return m.GetGauge().GetValue() }); err != nil {
		return nil, err
	}
	if err := writeMetric(m.SimulatedUsers, "simulated_users", func(m *dto.Metric) float64 { return m.GetGauge().GetValue() }); err != nil {
		return nil, err
	}

	// Get Counter values
	if err := writeMetric(m.PostsCreated, "posts", func(m *dto.Metric) float64 { return m.GetCounter().GetValue() }); err != nil {
		return nil, err
	}
	if err := writeMetric(m.CommentsCreated, "comments", func(m *dto.Metric) float64 { return m.GetCounter().GetValue() }); err != nil {
		return nil, err
	}
	if err := writeMetric(m.VotesRecorded, "votes", func(m *dto.Metric) float64 { return m.GetCounter().GetValue() }); err != nil {
		return nil, err
	}
	if err := writeMetric(m.TotalUsers, "total_users", func(m *dto.Metric) float64 { return m.GetCounter().GetValue() }); err != nil {
		return nil, err
	}

	// Get other gauge values
	if err := writeMetric(m.AverageResponseTime, "request_duration", func(m *dto.Metric) float64 { return m.GetGauge().GetValue() }); err != nil {
		return nil, err
	}
	if err := writeMetric(m.ErrorRate, "error_rate", func(m *dto.Metric) float64 { return m.GetGauge().GetValue() }); err != nil {
		return nil, err
	}

	return metrics, nil
}
