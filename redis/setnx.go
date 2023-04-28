package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ Lock = (*LockWithSetnx)(nil)

type LockWithSetnx struct {
	client  *redis.Client
	key     string
	expires time.Duration
}

func NewLockWithSetnx(client *redis.Client, key string, expires time.Duration) *LockWithSetnx {
	return &LockWithSetnx{
		client:  client,
		key:     key,
		expires: expires,
	}
}

func (l *LockWithSetnx) Lock() (bool, error) {
	ctx := context.Background()
	ok, err := l.client.SetNX(ctx, l.key, 1, l.expires).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (l *LockWithSetnx) Unlock() error {
	ctx := context.Background()
	_, err := l.client.Del(ctx, l.key).Result()
	return err
}
