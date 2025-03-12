package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	Key    string
	Value  string
	Ttl    time.Duration
	client *redis.Client
}

func NewRedisLock(client *redis.Client, key, value string, ttl time.Duration) *RedisLock {
	return &RedisLock{
		Key:    key,
		Value:  value,
		Ttl:    ttl,
		client: client,
	}
}

func (r *RedisLock) Aquire(ctx context.Context) (bool, error) {
	lockKey := fmt.Sprintf("lock:%v", r.Key)
	return r.client.SetNX(ctx, lockKey, r.Value, time.Duration(r.Ttl)*time.Second).Result()
}

func (r *RedisLock) Release(ctx context.Context) error {
	lockKey := fmt.Sprintf("lock:%v", r.Key)
	return r.client.Del(ctx, lockKey).Err()
}
