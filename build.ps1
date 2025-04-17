# Build script for microservices system
docker build -t inventory-service -f microservices-system/inventory/Dockerfile .
docker build -t order-service -f microservices-system/order/Dockerfile .