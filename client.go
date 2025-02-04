package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"storehouse/microservices-system/order/orderservice"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	client := orderservice.NewOrderServiceClient(conn)
	response, err := client.GetResponse(context.Background(), &orderservice.OrderResponse{
		Status:  "SUCCESS",
		Message: "Order placed successfully",
	})
	if err != nil {
		log.Fatalf("Error calling PlaceOrder: %v", err)
	}
	//response, err := client.PlaceOrder(context.Background(), &orderservice.OrderRequest{
	//	OrderId:    "123",
	//	CustomerId: "cust01",
	//	Items: []*orderservice.OrderItem{
	//		{ProductId: "prod01", Quantity: 2},
	//	},
	//})
	//if err != nil {
	//	log.Fatalf("Error calling PlaceOrder: %v", err)
	//}

	log.Printf("Order Response: Status=%s, Message=%s", response.Status, response.Message)
}
