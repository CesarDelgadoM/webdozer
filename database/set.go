package database

import (
	"context"
	"log"
)

type Set struct {
	*RedisPool
}

func NewSet(redis *RedisPool) *Set {
	return &Set{
		redis,
	}
}

func (s *Set) Add(key, value string) error {
	return s.client.SAdd(context.Background(), key, value).Err()
}

func (s *Set) Get(key, value string) error {
	return nil
}

func (s *Set) Exist(key, value string) bool {
	r, err := s.client.SIsMember(context.Background(), key, value).Result()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return r
}

func (s *Set) Del(key string) error {
	err := s.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}
