package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	idpb "banking-system/proto/idempotency"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache() *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return &Cache{rdb: rdb}
}

func (c *Cache) Check(userID int64, txID string) (*idpb.CheckTransactionResponse, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d:transaction:%s", userID, txID)
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return &idpb.CheckTransactionResponse{Exists: false, CachedBalance: 0}, nil
	}
	if err != nil {
		return nil, err
	}
	bal, _ := strconv.ParseFloat(val, 64)
	return &idpb.CheckTransactionResponse{Exists: true, CachedBalance: bal}, nil
}

func (c *Cache) Save(userID int64, txID string, balance float64) (*idpb.SaveTransactionResponse, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%d:transaction:%s", userID, txID)
	val := strconv.FormatFloat(balance, 'f', -1, 64)
	err := c.rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		return &idpb.SaveTransactionResponse{Success: false}, err
	}
	return &idpb.SaveTransactionResponse{Success: true}, nil
}
