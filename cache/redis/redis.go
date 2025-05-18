package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
	ctx context.Context
}

type IRedis interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration int) error
}

func NewRedis(rdb *redis.Client, ctx context.Context) IRedis {
	return &Redis{
		rdb: rdb,
		ctx: ctx,
	}
}

func (r *Redis) Get(key string) (string, error) {
	result, err := r.rdb.Get(r.ctx, key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (r *Redis) Set(key string, value string, expiration int) error {
	err := r.rdb.Set(r.ctx, key, value, time.Duration(expiration)*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}
