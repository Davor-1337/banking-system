package main

import (
	pb "banking-system/proto/db"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type server struct {
	pb.UnimplementedDBServiceServer
	db *Database
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	userID, err := s.db.GetUserByCredentials(req.Username, req.Password)
	if err != nil {
		return &pb.LoginResponse{Error: 1, Token: ""}, err
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour)

	if err := s.db.SaveToken(token, userID, expiresAt); err != nil {
		return &pb.LoginResponse{Error: 2, Token: ""}, err
	}

	return &pb.LoginResponse{Error: 0, Token: token, UserId: int32(userID)}, nil
}

func (s *server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, err := s.db.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Error: 1, UserId: 0}, err
	}

	return &pb.ValidateTokenResponse{Error: 0, UserId: int32(userID)}, nil
}

func (s *server) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {

	validateResp, err := s.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: req.Token})
	if err != nil || validateResp.Error != 0 {
		return &pb.DepositResponse{Error: 1, Balance: 0}, fmt.Errorf("invalid token")
	}
	userID := int(validateResp.UserId)

	if req.TransactionId == "" {
		return &pb.DepositResponse{Error: 5, Balance: 0}, fmt.Errorf("transaction_id required")
	}

	if err := s.db.UpdateBalance(userID, req.Amount); err != nil {
		return &pb.DepositResponse{Error: 3, Balance: 0}, err
	}

	newBalance, _ := s.db.GetBalance(userID)

	s.db.SaveTransaction(userID, req.TransactionId, "deposit", req.Amount, newBalance, req.Timestamp)

	return &pb.DepositResponse{Error: 0, Balance: newBalance}, nil
}

func (s *server) Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {

	validateResp, err := s.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: req.Token})
	if err != nil || validateResp.Error != 0 {
		return &pb.WithdrawResponse{Error: 1, Balance: 0}, fmt.Errorf("invalid token")
	}
	userID := int(validateResp.UserId)

	balance, err := s.db.GetBalance(userID)
	if err != nil {
		return &pb.WithdrawResponse{Error: 3, Balance: 0}, err
	}
	if balance < req.Amount {
		return &pb.WithdrawResponse{Error: 4, Balance: balance}, fmt.Errorf("insufficient funds")
	}

	if err := s.db.UpdateBalance(userID, -req.Amount); err != nil {
		return &pb.WithdrawResponse{Error: 5, Balance: 0}, err
	}

	newBalance, _ := s.db.GetBalance(userID)

	s.db.SaveTransaction(userID, req.TransactionId, "withdraw", req.Amount, newBalance, req.Timestamp)

	return &pb.WithdrawResponse{Error: 0, Balance: newBalance}, nil
}

func (s *server) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {

	userID, err := s.db.ValidateToken(req.Token)
	if err != nil {
		return &pb.GetBalanceResponse{Error: 1, Balance: 0}, err
	}

	balance, err := s.db.GetBalance(userID)
	if err != nil {
		return &pb.GetBalanceResponse{Error: 2, Balance: 0}, err
	}

	return &pb.GetBalanceResponse{Error: 0, Balance: balance}, nil
}
