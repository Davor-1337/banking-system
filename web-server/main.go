package main

import (
	dbpb "banking-system/proto/db"
	idpb "banking-system/proto/idempotency"
	"banking-system/web-server/handlers"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var dbClient dbpb.DBServiceClient
var idempotencyClient idpb.IdempotencyServiceClient

func main() {

	dbServiceHost := os.Getenv("DB_SERVICE_HOST")
	if dbServiceHost == "" {
		dbServiceHost = "localhost:50051"
	}

	dbConn, err := grpc.NewClient(dbServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to DB Service: %v", err)
	}
	defer dbConn.Close()
	dbClient = dbpb.NewDBServiceClient(dbConn)
	log.Printf("Connected to DB Service at %s", dbServiceHost)

	idempotencyServiceHost := os.Getenv("IDEMPOTENCY_SERVICE_HOST")
	if idempotencyServiceHost == "" {
		idempotencyServiceHost = "localhost:50052"
	}

	idConn, err := grpc.NewClient(idempotencyServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Idempotency Service: %v", err)
	}
	defer idConn.Close()
	idempotencyClient = idpb.NewIdempotencyServiceClient(idConn)
	log.Printf("Connected to Idempotency Service at %s", idempotencyServiceHost)

	// http routes
	loginHandler := &handlers.LoginHandler{DBClient: dbClient}
	http.Handle("/login", loginHandler)

	depositHandler := &handlers.DepositHandler{
		DBClient:          dbClient,
		IdempotencyClient: idempotencyClient,
	}
	http.Handle("/deposit", depositHandler)

	withdrawHandler := &handlers.WithdrawHandler{
		DBClient:          dbClient,
		IdempotencyClient: idempotencyClient,
	}
	http.Handle("/withdraw", withdrawHandler)

	log.Println("Web Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
