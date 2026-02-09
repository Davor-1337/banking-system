package handlers

import (
	dbpb "banking-system/proto/db"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type LoginHandler struct {
	DBClient dbpb.DBServiceClient
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Login request: username=%s", req.Username)

	resp, err := h.DBClient.Login(context.Background(), &dbpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Login failed: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": 1,
			"token": "",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": resp.Error,
		"token": resp.Token,
	})
}
