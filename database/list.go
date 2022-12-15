package database

import "context"

type List struct {
	*RedisPool
}

func NewList(redis *RedisPool) *List {
	return &List{
		redis,
	}
}

func (l *List) Add(key, value string) error {
	return l.client.LPush(context.Background(), key, value).Err()
}

func (l *List) Get(key, value string) error {
	return nil
}

func (l *List) Exist(key, value string) bool {
	return false
}

func (l *List) Del(key string) error {
	return nil
}
