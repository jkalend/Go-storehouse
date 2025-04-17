// Package server implements a gRPC server for the inventory service.
// It provides functionality for managing product inventory, including
// listing, retrieving details, modifying stock levels, and checking availability.
package server

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"storehouse/microservices-system/inventory/inventoryservice"
)

// InventoryServer implements the InventoryServiceServer gRPC interface.
// It manages product inventory through a PostgreSQL database connection.
type InventoryServer struct {
	inventoryservice.InventoryServiceServer
	logger *log.Logger // Logger for recording server operations
	conn   *pgx.Conn   // PostgreSQL database connection
}

// ListInventory returns a list of all products in the inventory.
// It retrieves the product ID, name, stock quantity, and price for each product.
//
// Parameters:
//   - ctx: The context for the request
//   - req: The empty request message
//
// Returns:
//   - InventoryListResponse containing all products' information
//   - Error if database operations fail
func (i *InventoryServer) ListInventory(context.Context, *inventoryservice.InventoryListRequest) (*inventoryservice.InventoryListResponse, error) {
	i.logger.Println("Listing inventory")

	var count int32
	row := i.conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM products")
	err := row.Scan(&count)
	if err != nil {
		return nil, err
	}
	i.logger.Printf("Found %d products\n", count)

	rows, err := i.conn.Query(context.Background(), "SELECT id, name, stock, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &inventoryservice.InventoryListResponse{}
	for rows.Next() {
		var id int32
		var name string
		var quantity int32
		var price float32
		err := rows.Scan(&id, &name, &quantity, &price)
		if err != nil {
			return nil, err
		}
		if i.logger.Prefix() == "DEBUG" {
			i.logger.Printf("Found product ID %d\n", id)
		}
		response.Inventory = append(response.Inventory, &inventoryservice.InventoryCheckResponse{
			ProductId: id,
			Name:      name,
			Quantity:  quantity,
			Price:     price,
		})
	}

	return response, nil
}

// GetDetails retrieves detailed information about a specific product.
//
// Parameters:
//   - ctx: The context for the request (unused but required by gRPC interface)
//   - in: Request containing the product ID to retrieve
//
// Returns:
//   - InventoryGetDetailResponse with complete product details
//   - Error if database operations fail or product not found
func (i *InventoryServer) GetDetails(_ context.Context, in *inventoryservice.InventoryGetDetailRequest) (*inventoryservice.InventoryGetDetailResponse, error) {
	i.logger.Printf("Showing details for product ID %d\n", in.ProductId)
	rows, err := i.conn.Query(context.Background(), "SELECT * FROM products WHERE id = $1", in.ProductId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var response *inventoryservice.InventoryGetDetailResponse
	for rows.Next() {
		var id int32
		var name string
		var description string
		var quantity int32
		var price float32
		err := rows.Scan(&id, &name, &description, &quantity, &price)
		if err != nil {
			return nil, err
		}
		if i.logger.Prefix() == "DEBUG" {
			i.logger.Printf("Showing details for ID %d\n", id)
		}
		response = &inventoryservice.InventoryGetDetailResponse{
			ProductId:   id,
			Name:        name,
			Description: description,
			Quantity:    quantity,
		}
	}
	return response, nil
}

// ModifyInventory updates the stock quantity for a specific product.
//
// Parameters:
//   - ctx: The context for the request (unused but required by gRPC interface)
//   - in: Request containing the product ID and new quantity
//
// Returns:
//   - InventoryModifyResponse with status information
//   - Error if database operations fail
func (i *InventoryServer) ModifyInventory(_ context.Context, in *inventoryservice.InventoryModifyRequest) (*inventoryservice.InventoryModifyResponse, error) {
	i.logger.Printf("Modifying inventory for product ID %d", in.ProductId)
	_, err := i.conn.Query(context.Background(), "UPDATE products SET stock = $1 WHERE id = $2", in.Quantity, in.ProductId)
	if err != nil {
		return nil, err
	}
	i.logger.Printf("Inventory updated for product ID %d\n", in.ProductId)
	return &inventoryservice.InventoryModifyResponse{
		Status:  "SUCCESS",
		Message: "Inventory updated successfully",
	}, nil
}

// CreateInventory adds a new product to the inventory.
//
// Parameters:
//   - ctx: The context for the request (unused but required by gRPC interface)
//   - in: Request containing product details (name, quantity, price)
//
// Returns:
//   - InventoryCreateResponse with the new product ID and status
//   - Error if database operations fail
func (i *InventoryServer) CreateInventory(_ context.Context, in *inventoryservice.InventoryCreateRequest) (*inventoryservice.InventoryCreateResponse, error) {
	i.logger.Println("Creating inventory")
	row, err := i.conn.Query(context.Background(), "INSERT INTO products (name, stock, price) VALUES ($1, $2, $3)", in.Name, in.Quantity, in.Price)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var id int32
	if row.Next() {
		err := row.Scan(&id)
		if err != nil {
			return nil, err
		}
		i.logger.Printf("Inventory created with ID %d\n", id)
	}

	return &inventoryservice.InventoryCreateResponse{
		ProductId: id,
		Status:    "SUCCESS",
		Message:   "Inventory created successfully",
	}, nil
}

// DeleteInventory removes a product from the inventory.
//
// Parameters:
//   - ctx: The context for the request (unused but required by gRPC interface)
//   - in: Request containing the product ID to delete
//
// Returns:
//   - InventoryDeleteResponse with status information
//   - Error if database operations fail
func (i *InventoryServer) DeleteInventory(_ context.Context, in *inventoryservice.InventoryDeleteRequest) (*inventoryservice.InventoryDeleteResponse, error) {
	i.logger.Printf("Deleting inventory for product ID %d\n", in.ProductId)
	_, err := i.conn.Query(context.Background(), "DELETE FROM products WHERE id = $1", in.ProductId)
	if err != nil {
		return nil, err
	}
	if i.logger.Prefix() == "DEBUG" {
		i.logger.Printf("Inventory deleted for product ID %d\n", in.ProductId)
	}
	return &inventoryservice.InventoryDeleteResponse{
		Status:  "SUCCESS",
		Message: "Inventory deleted successfully",
	}, nil
}

// NewInventoryServer creates and returns a new instance of InventoryServer.
//
// Parameters:
//   - logger: Logger for recording server operations
//   - conn: PostgreSQL database connection
//
// Returns:
//   - Pointer to a new InventoryServer instance
func NewInventoryServer(logger *log.Logger, conn *pgx.Conn) *InventoryServer {
	return &InventoryServer{logger: logger, conn: conn}
}

// CheckStock verifies if a product has enough stock to fulfill a requested quantity.
//
// Parameters:
//   - ctx: The context for the request (unused but required by gRPC interface)
//   - in: Request containing the product ID and quantity to check
//
// Returns:
//   - InventoryStockResponse with availability status and current price
//   - Error if database operations fail or product not found
func (i *InventoryServer) CheckStock(_ context.Context, in *inventoryservice.InventoryStockRequest) (*inventoryservice.InventoryStockResponse, error) {
	i.logger.Printf("Checking inventory for product ID %d\n", in.ProductId)
	rows, err := i.conn.Query(context.Background(), "SELECT stock, price FROM products WHERE id = $1", in.ProductId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quantity int32
	var price float32
	if rows.Next() {
		err := rows.Scan(&quantity, &price)
		if err != nil {
			return nil, err
		}
	}
	i.logger.Printf("Found %d items in stock for product ID %d\n", quantity, in.ProductId)
	if i.logger.Prefix() == "DEBUG" {
		i.logger.Printf("Found %d items in stock for product ID %d\n", quantity, in.ProductId)
	}
	return &inventoryservice.InventoryStockResponse{
		InStock: quantity >= in.Quantity,
		Price:   price,
	}, nil
}
