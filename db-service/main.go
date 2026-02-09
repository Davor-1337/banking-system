package main

import (
	pb "banking-system/proto/db"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

const port = ":50051"

func main() {

	dsn := "bankuser:bankpass@tcp(mysql:3306)/banking"
	database, err := NewDatabase(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer database.Close()

	log.Println("Successfully connected to MySQL")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDBServiceServer(grpcServer, &server{db: database})

	log.Printf("DB Service listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
