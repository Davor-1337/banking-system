package main

import (
	"log"
	"net"

	idpb "banking-system/proto/idempotency"

	"google.golang.org/grpc"
)

const port = ":50052"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	idpb.RegisterIdempotencyServiceServer(grpcServer, NewIdempotencyServer())

	log.Printf("Idempotency Service listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
