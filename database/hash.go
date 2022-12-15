package database

import (
	"context"
	"log"
)

type Hash struct {
	name string
	ctx  context.Context
	*RedisPool
}

func NewHash(name string, client *RedisPool) Hash {

	return Hash{
		name,
		context.Background(),
		client,
	}
}

func (h *Hash) Set(key, value string) {

	_, err := h.client.HSet(h.ctx, h.name, key, value).Result()
	if err != nil {
		log.Println("Failed operation Set:", err)
	}
}
