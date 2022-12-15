package database

import (
	"log"

	"github.com/go-redis/redis/v8"
)

type Repository interface {
	Add(key, value string) error
	Get(key, value string) error
	Exist(key, value string) bool
	Del(key string) error
}

type RedisPool struct {
	client *redis.Client
}

func NewRedisPool(addr string, procs int) *RedisPool {
	log.Println("Max number of connections for redis:", procs)

	pool := redis.NewClient(&redis.Options{
		Addr:     addr,
		PoolSize: procs,
	})

	return &RedisPool{
		pool,
	}
}

func (r *RedisPool) Close() {
	err := r.client.Close()
	if err != nil {
		log.Fatal("Failed to close redis connection")
	}
}
