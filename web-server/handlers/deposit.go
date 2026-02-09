package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	dbpb "banking-system/proto/db"
	idpb "banking-system/proto/idempotency"
)

type DepositHandler struct {
	DBClient          dbpb.DBServiceClient
	IdempotencyClient idpb.IdempotencyServiceClient
}

func (h *DepositHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Id        string  `json:"id"`
		Token     string  `json:"token"`
		Amount    float64 `json:"amount"`
		Timestamp int64   `json:"timestamp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Deposit request: transaction_id=%s, amount=%.2f", req.Id, req.Amount)

	validateResp, err := h.DBClient.ValidateToken(context.Background(), &dbpb.ValidateTokenRequest{
		Token: req.Token,
	})
	if err != nil || validateResp.Error != 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"error": 1, "balance": 0})
		return
	}
	userID := validateResp.UserId

	idempotencyResp, err := h.IdempotencyClient.CheckTransaction(context.Background(), &idpb.CheckTransactionRequest{
		UserId:        userID,
		TransactionId: req.Id,
	})
	if err == nil && idempotencyResp.Exists {
		json.NewEncoder(w).Encode(map[string]interface{}{"error": 0, "balance": idempotencyResp.CachedBalance})
		return
	}

	depositResp, err := h.DBClient.Deposit(context.Background(), &dbpb.DepositRequest{
		Token:         req.Token,
		TransactionId: req.Id,
		Amount:        req.Amount,
		Timestamp:     req.Timestamp,
	})
	if err != nil || depositResp.Error != 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"error": 1, "balance": 0})
		return
	}

	_, _ = h.IdempotencyClient.SaveTransaction(context.Background(), &idpb.SaveTransactionRequest{
		UserId:        userID,
		TransactionId: req.Id,
		Balance:       depositResp.Balance,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{"error": 0, "balance": depositResp.Balance})
}
