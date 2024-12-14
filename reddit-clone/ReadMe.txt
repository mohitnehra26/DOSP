Reddit Clone with Proto.Actor
A distributed Reddit clone implementation using Proto.Actor framework in Go.
Prerequisites
Go 1.23.3 or later
Protocol Buffers compiler (protoc)
Git
Installation
Clone the repository:
bash
git clone <your-repository-url>
cd reddit-clone

Install Go dependencies:
bash
go mod download
go get github.com/asynkron/protoactor-go
go get github.com/prometheus/client_golang/prometheus

Generate Protocol Buffer files:
bash
protoc --go_out=. --go-grpc_out=. api/proto/*.proto

Project Structure
text
reddit-clone/
├── api/
│   └── proto/          # Protocol Buffer definitions
├── cmd/
│   ├── engine/         # Engine service entry point
│   └── simulator/      # Client simulator entry point
├── internal/
│   ├── actor/          # Actor implementations
│   ├── common/         # Shared types and utilities
│   └── store/          # Data storage implementation
└── pkg/
    ├── metrics/        # Prometheus metrics
    └── utils/          # Utility functions

Running the Project
Start the Engine:
bash
go run cmd/engine/main.go

Start the Simulator in a new terminal:
bash
go run cmd/simulator/main.go

Monitoring
Access metrics through Prometheus endpoints:
Engine metrics: http://localhost:2112/metrics
Simulator metrics: http://localhost:2113/metrics
System Requirements
Based on metrics:
Memory: ~45MB RAM (Engine: 26.7MB, Simulator: 24.2MB)
CPU: Minimal (0.015625 seconds total CPU time)
Threads: 16-20 threads
Goroutines: 23-27 concurrent goroutines