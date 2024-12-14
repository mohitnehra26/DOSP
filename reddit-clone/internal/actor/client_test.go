package actor

import (
	"github.com/asynkron/protoactor-go/actor"

	"reddit-clone/internal/common"
	"reddit-clone/pkg/metrics"
	"testing"
	"time"
)

func TestClientActorReceive(t *testing.T) {
	system := actor.NewActorSystem()
	metrics := metrics.NewRedditMetrics()
	behavior := &common.ClientBehavior{
		PostProbability:    0.3,
		CommentProbability: 0.3,
		VoteProbability:    0.3,
		JoinProbability:    0.1,
		ActiveHours:        []int{9, 10, 11, 12, 13, 14, 15, 16, 17},
		Persona:            "Casual",
	}

	enginePID := system.Root.Spawn(actor.PropsFromFunc(func(context actor.Context) {}))

	props := actor.PropsFromProducer(func() actor.Actor {
		return NewClientActor("user1", "testuser", enginePID, behavior, metrics)
	})
	pid := system.Root.Spawn(props)

	system.Root.Send(pid, &common.SimulateAction{Timestamp: time.Now()})

	// Wait for a short time to allow the message to be processed
	time.Sleep(100 * time.Millisecond)

	essentialMetrics, err := metrics.GetEssentialMetrics()
	if err != nil {
		// Handle the error appropriately, e.g.:
		t.Fatalf("Failed to get essential metrics: %v", err)
	}
	if essentialMetrics["active_users"] != 1 {
		t.Error("Expected active users to be 1")
	}
}
