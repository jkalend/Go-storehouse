package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"storehouse/microservices-system/inventory/inventoryservice"
	"storehouse/microservices-system/order/orderservice"
	"storehouse/microservices-system/payment/paymentservice"
)

type OrderServer struct {
	orderservice.OrderServiceServer
	logger *log.Logger
}

func NewOrderServer(logger *log.Logger) *OrderServer {
	return &OrderServer{logger: logger}
}

func (s *OrderServer) PlaceOrder(ctx context.Context, req *orderservice.OrderRequest) (*orderservice.OrderResponse, error) {
	// Check inventory
	inventoryConn, err := grpc.NewClient("localhost:50052")
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
	paymentConn, err := grpc.NewClient("localhost:50053")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	}
	defer func(paymentConn *grpc.ClientConn) {
		err := paymentConn.Close()
		if err != nil {
			s.logger.Printf("failed to close payment connection: %v", err)
		}
	}(paymentConn)

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

func (s *OrderServer) GetResponse(ctx context.Context, res *orderservice.OrderResponse) (*orderservice.OrderResponse, error) {
	s.logger.Printf("Received response: Status=%s, Message=%s", res.Status, res.Message)
	return &orderservice.OrderResponse{Status: res.GetStatus(), Message: res.GetMessage()}, nil
}
