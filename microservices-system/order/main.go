package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/url"
	"os"
	"storehouse/microservices-system/order/orderservice"
	"storehouse/microservices-system/order/server"
)

func main() {
	logger := log.New(log.Writer(), "order-service ", log.LstdFlags)
	logger.Println("Starting order service")
	srv := grpc.NewServer()

	var connString string
	if os.Getenv("POSTGRESQL_USER") == "" {
		panic("POSTGRESQL_USER not set")
	}
	if os.Getenv("POSTGRESQL_PASSWORD") == "" {
		panic("POSTGRESQL_PASSWORD not set")
	}

	username := url.QueryEscape(os.Getenv("POSTGRESQL_USER"))
	password := url.QueryEscape(os.Getenv("POSTGRESQL_PASSWORD"))

	logger.Println("Connecting to database")
	connString = "postgres://" + username + ":" + password + "@storehouse-inventory:5432/storehouse_inventory"
	logger.Printf("Connecting to database with connection string: %s\n", connString)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())

	logger.Println("Registering order service")
	orderservice.RegisterOrderServiceServer(srv, server.NewOrderServer(logger, conn))

	logger.Println("Listening on port 50051")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}

	defer func() {
		logger.Println("Stopping order service")
		srv.Stop()
	}()
}
