# This script builds the Docker images for the inventory and order services, creates a Docker network, and runs the services in containers.
# The main config is for k8s, this is for testing locally
./build.ps1
docker network create microservices
docker stop inventory-service
docker stop order-service
docker run -d --rm --network microservices --name inventory-service -p 50052:50052 inventory-service:latest
docker run -d --rm --network microservices --name order-service -p 50051:50051 order-service:latest
