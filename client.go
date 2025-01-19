package main

import (
	"context"
	"first/microservices-system/order/orderservice"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := orderservice.NewOrderServiceClient(conn)
	response, err := client.PlaceOrder(context.Background(), &orderservice.OrderRequest{
		OrderId:    "123",
		CustomerId: "cust01",
		Items: []*orderservice.OrderItem{
			{ProductId: "prod01", Quantity: 2},
		},
	})
	if err != nil {
		log.Fatalf("Error calling PlaceOrder: %v", err)
	}

	log.Printf("Order Response: Status=%s, Message=%s", response.Status, response.Message)
}
