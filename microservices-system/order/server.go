package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"first/microservices-system/inventory/inventoryservice"
	"first/microservices-system/order/orderservice"
	"first/microservices-system/payment/paymentservice"
	"google.golang.org/grpc"
)

type orderServer struct {
	orderservice.OrderServiceServer
}

func (s *orderServer) PlaceOrder(ctx context.Context, req *orderservice.OrderRequest) (*orderservice.OrderResponse, error) {
	// Check inventory
	inventoryConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %v", err)
	}
	defer inventoryConn.Close()

	inventoryClient := inventoryservice.NewInventoryServiceClient(inventoryConn)
	inventoryRes, err := inventoryClient.CheckInventory(ctx, &inventoryservice.InventoryCheckRequest{
		ProductId: req.Items[0].ProductId,
		Quantity:  req.Items[0].Quantity,
	})
	if err != nil || !inventoryRes.InStock {
		return &orderservice.OrderResponse{Status: "FAILED", Message: "Insufficient stock"}, nil
	}

	// Process payment
	paymentConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	paymentClient := paymentservice.NewPaymentServiceClient(paymentConn)
	paymentRes, err := paymentClient.ProcessPayment(ctx, &paymentservice.PaymentRequest{
		OrderId: req.OrderId,
		Amount:  100.0, // Example amount
	})
	if err != nil || paymentRes.Status != "SUCCESS" {
		return &orderservice.OrderResponse{Status: "FAILED", Message: "Payment failed"}, nil
	}

	return &orderservice.OrderResponse{Status: "SUCCESS", Message: "Order placed successfully"}, nil
}

func main() {
	server := grpc.NewServer()
	orderservice.RegisterOrderServiceServer(server, &orderServer{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Order Service is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
