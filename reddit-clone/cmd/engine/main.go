package main

import (
	"context"
	"fmt"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	internalActor "reddit-clone/internal/actor" // Alias the import
	"reddit-clone/internal/store/memory"
	"reddit-clone/pkg/metrics"
)

func main() {
	// Initialize metrics
	metricsCollector := metrics.NewRedditMetrics()

	// Start metrics endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Initialize actor system
	system := actor.NewActorSystem()

	// Configure remote
	remoteConfig := remote.Configure("127.0.0.1", 8090)

	// Create remote
	remoting := remote.NewRemote(system, remoteConfig)

	// Create new engine actor
	engineActor := internalActor.NewEngineActor(
		memory.NewMemoryStore(),
		metricsCollector,
	)

	// Create props
	props := actor.PropsFromProducer(func() actor.Actor {
		return engineActor
	})

	remoting.Register("engine", props)
	remoting.Start()

	// Spawn the engine actor
	pid, err := system.Root.SpawnNamed(props, "engine")
	if err != nil {
		log.Fatalf("Failed to spawn engine actor: %v", err)
	}

	fmt.Printf("Engine started. PID: %v\n", pid)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down the engine...")
			system.Root.Stop(pid)
			return
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)

		cancel() // Signal all components to shut down

		fmt.Println("Stopping remoting...")
		remoting.Shutdown(true)
		
		fmt.Println("Shutting down actor system...")
		system.Shutdown()

		fmt.Println("Shutdown complete.")
		os.Exit(0)
	}()

	// Keep the process running until shutdown is triggered
	select {}
}
