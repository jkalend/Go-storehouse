protoc --go_out=. --go-grpc_out=. order_service.proto
protoc --go_out=. --go-grpc_out=. inventory_service.proto
protoc --go_out=. --go-grpc_out=. payment_service.proto
Write-Host "Protobuf files generated successfully! Yipee!"