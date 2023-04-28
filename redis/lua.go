package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ Lock = (*LockWithLua)(nil)

type LockWithLua struct {
	client  *redis.Client
	key     string
	expires time.Duration
}

func (l LockWithLua) Lock() (bool, error) {
	ctx := context.Background()
	result, err := l.client.Eval(ctx, `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
            return "OK"
        elseif redis.call("SETNX", KEYS[1], ARGV[1]) == 1 then
            redis.call("PEXPIRE", KEYS[1], ARGV[2])
            return "OK"
        else
            return nil
        end
	`, []string{l.key}, 1, l.key, 1, l.expires).Result()
	if err != nil {
		return false, err
	}
	return result != nil, nil
}

func (l LockWithLua) Unlock() error {

}
