package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"storehouse/microservices-system/order/orderservice"
	"storehouse/microservices-system/order/server"
)

func main() {
	logger := log.New(log.Writer(), "order-service", log.LstdFlags)
	srv := grpc.NewServer()
	orderservice.RegisterOrderServiceServer(srv, server.NewOrderServer(logger))
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
