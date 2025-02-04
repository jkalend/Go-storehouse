package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log"
	"storehouse/microservices-system/inventory/inventoryservice"
	"storehouse/microservices-system/order/orderservice"
	"storehouse/microservices-system/payment/paymentservice"
)

type InventoryServer struct {
	inventoryservice.InventoryServiceServer
	logger *log.Logger
}

func NewInventoryServer(logger *log.Logger) *InventoryServer {
	return &InventoryServer{logger: logger}
}
