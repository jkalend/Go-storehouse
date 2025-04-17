// Package main provides a client application for demonstrating and testing the
// storehouse microservices system. It includes examples of connecting to and
// interacting with both the order and inventory services.
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"storehouse/microservices-system/inventory/inventoryservice"
	"storehouse/microservices-system/order/orderservice"
)

func main() {
	// Connect to the Order Service
	conn, err := grpc.NewClient("localhost:30051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close order service connection: %v", err)
		}
	}(conn)

	// Place an order with multiple items
	client := orderservice.NewOrderServiceClient(conn)
	response, err := client.PlaceOrder(context.Background(), &orderservice.OrderRequest{
		CustomerId: 1,
		Items: []*orderservice.OrderItem{
			{ProductId: 1, Quantity: 2},
			{ProductId: 2, Quantity: 5},
		},
	})
	if err != nil {
		log.Fatalf("Error calling PlaceOrder: %v", err)
	}
	fmt.Println(response.Message)
	fmt.Println(response.OrderId)

	// Connect to the Inventory Service
	conn2, err := grpc.NewClient("localhost:30052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func(conn2 *grpc.ClientConn) {
		err := conn2.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn2)

	// List all inventory items
	client2 := inventoryservice.NewInventoryServiceClient(conn2)
	response2, err := client2.ListInventory(context.Background(), &inventoryservice.InventoryListRequest{})
	if err != nil {
		log.Fatalf("Error calling ListInventory: %v", err)
	}

	// Display inventory items information
	for _, item := range response2.Inventory {
		log.Printf("Product ID: %d, Name: %s, Quantity: %d, Price: %f", item.ProductId, item.Name, item.Quantity, item.Price)
	}
}
