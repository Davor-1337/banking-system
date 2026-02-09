package main

import (
	idpb "banking-system/proto/idempotency"
	"context"
)

type idempotencyServer struct {
	idpb.UnimplementedIdempotencyServiceServer
	cache *Cache
}

func NewIdempotencyServer() *idempotencyServer {
	return &idempotencyServer{cache: NewCache()}
}

func (s *idempotencyServer) CheckTransaction(ctx context.Context, req *idpb.CheckTransactionRequest) (*idpb.CheckTransactionResponse, error) {
	return s.cache.Check(int64(req.UserId), req.TransactionId)
}

func (s *idempotencyServer) SaveTransaction(ctx context.Context, req *idpb.SaveTransactionRequest) (*idpb.SaveTransactionResponse, error) {
	return s.cache.Save(int64(req.UserId), req.TransactionId, req.Balance)
}
