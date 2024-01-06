generate:
	protoc -I ./proto \
      --go_out=paths=import:./stubs \
      --go_opt=paths=source_relative \
      --go-grpc_out=paths=import:./stubs \
      --go-grpc_opt=paths=source_relative \
      ./proto/booking-service/v1/booking.proto

build-server:
	go build -o bin/server ./server

build-client:
	go build -o bin/client ./client

run-test:
	go test -v ./server/internal/apis_test

docker-network:
	sudo docker network create book-seat-network

build-docker-server:
	sudo docker build -t book-my-seat-img -f server/Dockerfile .

run-docker-server:
	 sudo docker rm -f book-my-seat && sudo docker run -it --name book-my-seat -p 50052:50051 --network book-seat-network book-my-seat-img

build-docker-client:
	sudo docker build -t book-my-seat-client-img -f client/Dockerfile .

run-docker-client:
	sudo docker rm -f book-my-seat-client && sudo docker run -it --name book-my-seat-client --network book-seat-network book-my-seat-client-img
