# Go-storehouse

## Overview
Go-storehouse is a proof of concept (PoC) project designed to explore and demonstrate the implementation of microservices using Go, gRPC, and Kubernetes. This project simulates a storage/inventory management system with distributed services communicating via gRPC.

## Project Purpose
The primary goal of this project is educational - to gain hands-on experience with:
- Developing microservices in Go
- Implementing service communication using gRPC
- Containerizing applications
- Orchestrating containers with Kubernetes
- Designing resilient distributed systems

## Architecture
The project consists of several microservices that work together:
- **Inventory Service**: Manages product inventory and stock levels
- **Order Service**: Handles customer orders
- **Payment Service**: Processes payments
- **User Service**: Manages user accounts and authentication
- **API Gateway**: Provides a unified entry point for external clients

> **Note**: The payment service, user service and API gateway are missing as this is only a proof of concept.

## Technologies
- **Language**: Go
- **Communication**: gRPC for inter-service communication
- **Containerization**: Docker
- **Orchestration**: Kubernetes

## Getting Started

### Prerequisites
- Go 1.16+
- Docker and Docker Compose
- Kubernetes cluster or Minikube
- kubectl configured

### Setup
Set up for windows:

1. Clone the repository:
   ```
   git clone https://github.com/jkalend/Go-storehouse.git
   cd Go-storehouse
   ```

2. Generate protobuf files:
   ```
   .\gen_proto.ps1
   ```
   
3. Build Docker images:
   ```
    .\build.ps1
    ```
   
4. Start services on a running Kubernetes cluster:
   ```
   .\run-k8s.ps1
   ```

5. Compile and run the testing client:
   ```
   go run client.go
   ```

### Database Initialization
The system automatically initializes the database on startup:
- Database schema is created if it doesn't exist
- Sample data is seeded including test products and users
- The seeding logic checks if data already exists to avoid duplication

## Development Status
This project is a proof of concept and is not intended for production use. It demonstrates the architecture and implementation patterns for a microservices-based system.

### Limitations
- Payment processing is not implemented (intentionally omitted for the PoC)
- Limited error handling and recovery mechanisms
- Minimal security implementations, such as authentication and authorization or even secrets management

## Contributing
Contributions to extend this proof of concept are welcome. Please feel free to submit issues and pull requests.

## License
[MIT License](LICENSE)
