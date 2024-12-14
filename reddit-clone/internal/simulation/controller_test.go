package simulation

import (
	"github.com/asynkron/protoactor-go/actor"
	"reddit-clone/pkg/metrics"
	"testing"
	"time"
)

func TestSimulationControllerStart(t *testing.T) {
	system := actor.NewActorSystem()
	enginePID := system.Root.Spawn(actor.PropsFromFunc(func(context actor.Context) {}))
	metrics := metrics.NewRedditMetrics()

	controller := NewSimulationController(system, enginePID, metrics)

	err := controller.Start(10)
	if err != nil {
		t.Fatalf("Failed to start simulation: %v", err)
	}

	time.Sleep(1 * time.Second) // Give some time for the simulation to run

	if len(controller.clients) != 10 {
		t.Errorf("Expected 10 clients, got %d", len(controller.clients))
	}

	essentialMetrics, err := metrics.GetEssentialMetrics()
	if err != nil {
		// Handle the error appropriately, e.g.:
		t.Fatalf("Failed to get essential metrics: %v", err)
	}
	if essentialMetrics["simulated_users"] != 10 {
		t.Errorf("Expected 10 simulated users, got %f", essentialMetrics["simulated_users"])
	}
}
