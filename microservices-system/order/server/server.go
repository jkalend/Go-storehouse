// Package server implements the order service gRPC server functionality.
// It handles processing customer orders, checking inventory availability,
// and communicating with other microservices in the system.
package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"storehouse/microservices-system/inventory/inventoryservice"
	"storehouse/microservices-system/order/orderservice"
)

// OrderServer implements the OrderServiceServer interface and manages order processing.
// It maintains a connection to the database and handles order-related operations.
type OrderServer struct {
	orderservice.OrderServiceServer
	logger *log.Logger
	conn   *pgx.Conn
}

// NewOrderServer creates and returns a new OrderServer instance with the provided logger
// and database connection.
func NewOrderServer(logger *log.Logger, conn *pgx.Conn) *OrderServer {
	return &OrderServer{logger: logger, conn: conn}
}

// PlaceOrder processes an order request from a customer.
// It validates inventory availability, calculates the total amount,
// creates an order record in the database, and returns the result.
func (s *OrderServer) PlaceOrder(ctx context.Context, req *orderservice.OrderRequest) (*orderservice.OrderResponse, error) {
	// Log order processing start
	s.logger.Println("PlaceOrder")

	// Establish connection to the inventory service
	inventoryConn, err := grpc.NewClient("inventory-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %v", err)
	}
	defer func(inventoryConn *grpc.ClientConn) {
		err := inventoryConn.Close()
		if err != nil {

		}
	}(inventoryConn)

	s.logger.Println("Checking stock")
	inventoryClient := inventoryservice.NewInventoryServiceClient(inventoryConn)

	// Calculate total order amount and verify stock availability
	var amount float32 = 0.0

	for _, item := range req.Items {
		s.logger.Printf("Adding item: %v", item)
		inventoryRes, err := inventoryClient.CheckStock(ctx, &inventoryservice.InventoryStockRequest{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
		s.logger.Printf("Inventory response: %v", inventoryRes)
		s.logger.Printf("Inventory error: %v", err)
		if err != nil || !inventoryRes.InStock {
			// cancels whole order, could be improved to cancel only the item
			return &orderservice.OrderResponse{Status: "FAILED", Message: "Insufficient stock"}, nil
		}
		amount += inventoryRes.Price * float32(item.Quantity)
	}

	// Extract product IDs for database storage
	ids := make([]int32, len(req.Items))
	for i, item := range req.Items {
		ids[i] = item.ProductId
	}

	// Calculate total quantity across all items
	var quantity int32 = 0
	for _, item := range req.Items {
		quantity += item.Quantity
	}

	// Create order record in the database
	s.logger.Printf("Placing order for: CustomerId=%d", req.CustomerId)
	row := s.conn.QueryRow(ctx, "INSERT INTO orders (user_id, product_ids, quantity, total) VALUES ($1, $2, $3, $4) RETURNING id", req.CustomerId, ids, quantity, amount)

	var orderId int32
	if err := row.Scan(&orderId); err != nil {
		return &orderservice.OrderResponse{Status: "FAILED", Message: "Failed to place order in database"}, err
	}
	s.logger.Printf("Order %d placed for: CustomerId=%s", orderId, req.CustomerId)

	// Payment processing placeholder for future implementation
	//paymentConn, err := grpc.NewClient("localhost:50053")
	//if err != nil {
	//	return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	//}
	//defer func(paymentConn *grpc.ClientConn) {
	//	err := paymentConn.Close()
	//	if err != nil {
	//		s.logger.Printf("failed to close payment connection: %v", err)
	//	}
	//}(paymentConn)
	//
	//paymentClient := paymentservice.NewPaymentServiceClient(paymentConn)
	//paymentRes, err := paymentClient.ProcessPayment(ctx, &paymentservice.PaymentRequest{
	//	OrderId: req.OrderId,
	//	Amount:  amount,
	//})
	//if err != nil || paymentRes.Status != "SUCCESS" {
	//	return &orderservice.OrderResponse{Status: "FAILED", Message: "Payment failed"}, err
	//}

	return &orderservice.OrderResponse{OrderId: orderId, Status: "SUCCESS", Message: "Order placed successfully"}, nil
}

// GetResponse handles retrieving order response information.
// This is primarily used for debugging and monitoring the order service.
func (s *OrderServer) GetResponse(_ context.Context, res *orderservice.OrderResponse) (*orderservice.OrderResponse, error) {
	s.logger.Printf("Received response: Status=%s, Message=%s", res.Status, res.Message)
	return &orderservice.OrderResponse{Status: res.GetStatus(), Message: res.GetMessage()}, nil
}
