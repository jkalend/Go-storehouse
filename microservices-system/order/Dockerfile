# Use official Go image as base
FROM golang:1.23 AS builder

# Set working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application files
COPY . .

WORKDIR /app/microservices-system/order

# Build the Go application
RUN go build -o order-service .

# Use a lightweight image for the final container
FROM ubuntu:22.10

# Set working directory in container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/microservices-system/order/order-service .

# Expose gRPC port
EXPOSE 50051

# Command to run the application
CMD ["./order-service"]
