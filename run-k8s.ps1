# This script is used to deploy the microservices on a running Kubernetes cluster.
kubectl apply -f microservices-system/inventory/server/db/storehouse-inventory.yaml
kubectl apply -f microservices-system/inventory/inventory-service.yaml
kubectl apply -f microservices-system/order/order-service.yaml