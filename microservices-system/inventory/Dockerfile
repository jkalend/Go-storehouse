# Use official Go image as base
FROM golang:1.23 AS builder

# Set working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application files
COPY . .

WORKDIR /app/microservices-system/inventory

# Build the Go application
RUN go build -o inventory-service .

# Use a lightweight image for the final container
FROM ubuntu:22.10

# Set working directory in container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/microservices-system/inventory/inventory-service .

# copy the DB init files
COPY --from=builder /app/microservices-system/inventory/server/db server/db

ENV POSTGRESQL_USER='meowiiik'
ENV POSTGRESQL_PASSWORD='meowiiik#397Ang7Rt'

# Expose gRPC port
EXPOSE 50052

# Command to run the application
CMD ["./inventory-service"]
