generate:
	protoc -I ./proto \
      --go_out=paths=import:./stubs \
      --go_opt=paths=source_relative \
      --go-grpc_out=paths=import:./stubs \
      --go-grpc_opt=paths=source_relative \
      ./proto/booking-service/v1/booking.proto