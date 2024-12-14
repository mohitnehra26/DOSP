package metrics

import (
	"testing"
)

func TestUpdateActiveUsers(t *testing.T) {
	metrics := NewRedditMetrics()

	// Start with no active users
	metrics.UpdateActiveUsers(5)
	essentialMetrics, err := metrics.GetEssentialMetrics()
	if err != nil {
		t.Fatalf("Failed to get essential metrics: %v", err)
	}
	if essentialMetrics["active_users"] != 5 {
		t.Errorf("Expected active users to be 5, got %f", essentialMetrics["active_users"])
	}

	// Remove some users
	metrics.UpdateActiveUsers(-3)
	essentialMetrics, err = metrics.GetEssentialMetrics()
	if err != nil {
		t.Fatalf("Failed to get essential metrics: %v", err)
	}
	if essentialMetrics["active_users"] != 2 {
		t.Errorf("Expected active users to be 2, got %f", essentialMetrics["active_users"])
	}

	// Attempt to remove more users than available
	metrics.UpdateActiveUsers(-10)
	essentialMetrics, err = metrics.GetEssentialMetrics()
	if err != nil {
		t.Fatalf("Failed to get essential metrics: %v", err)
	}
	if essentialMetrics["active_users"] != 0 {
		t.Errorf("Expected active users to be 0, got %f", essentialMetrics["active_users"])
	}

	// Ensure Total Users only increases
	initialTotal := essentialMetrics["total_users"]
	metrics.UpdateActiveUsers(3)
	essentialMetrics, err = metrics.GetEssentialMetrics()
	if err != nil {
		t.Fatalf("Failed to get essential metrics: %v", err)
	}
	if essentialMetrics["total_users"] != initialTotal+3 {
		t.Errorf("Expected total users to increase by 3")
	}
}
func TestGetEssentialMetrics(t *testing.T) {
	m := NewRedditMetrics()

	m.ActiveUsers.Set(5)
	m.SimulatedUsers.Set(10)
	m.PostsCreated.Add(15)
	m.CommentsCreated.Add(20)
	m.VotesRecorded.Add(25)
	m.TotalUsers.Add(30)
	m.AverageResponseTime.Set(1.5)
	m.ErrorRate.Set(0.05)

	result, err := m.GetEssentialMetrics()
	if err != nil {
		t.Fatalf("Failed to get essential metrics: %v", err)
	}

	expectedMetrics := map[string]float64{
		"active_users":     5,
		"simulated_users":  10,
		"posts":            15,
		"comments":         20,
		"votes":            25,
		"total_users":      30,
		"request_duration": 1.5,
		"error_rate":       0.05,
	}

	for key, expected := range expectedMetrics {
		if result[key] != expected {
			t.Errorf("Expected %s to be %f, got %f", key, expected, result[key])
		}
	}
}
