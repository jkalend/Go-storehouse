package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"storehouse/microservices-system/inventory/inventoryservice"
	"storehouse/microservices-system/inventory/server"
)

// initDatabase initializes the database schema and seeds it with initial data
func initDatabase(ctx context.Context, conn *pgx.Conn, logger *log.Logger) error {
	logger.Println("Initializing database")

	// Read and execute schema SQL
	schemaPath := filepath.Join("server", "db", "Schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		logger.Printf("Error reading schema file: %v", err)
		return err
	}

	_, err = conn.Exec(ctx, string(schemaSQL))
	if err != nil {
		logger.Printf("Error executing schema SQL: %v", err)
		return err
	}
	logger.Println("Database schema created successfully")

	// Read and execute seeder SQL
	seederPath := filepath.Join("server", "db", "Seeder.sql")
	seederSQL, err := os.ReadFile(seederPath)
	if err != nil {
		logger.Printf("Error reading seeder file: %v", err)
		return err
	}

	_, err = conn.Exec(ctx, string(seederSQL))
	if err != nil {
		logger.Printf("Error executing seeder SQL: %v", err)
		return err
	}
	logger.Println("Database seeded successfully")

	return nil
}

func main() {
	logger := log.New(log.Writer(), "DEBUG", log.LstdFlags)
	srv := grpc.NewServer()
	defer func() {
		logger.Println("Stopping inventory service")
		srv.Stop()
	}()

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

	// Initialize the database
	if err := initDatabase(context.Background(), conn, logger); err != nil {
		logger.Printf("Warning: Database initialization failed: %v", err)
		// Continue execution even if initialization fails
	}

	logger.Println("Registering inventory service")
	inventoryservice.RegisterInventoryServiceServer(srv, server.NewInventoryServer(logger, conn))
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		panic(err)
	}
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
