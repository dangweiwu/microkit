package redisx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewClient(cfg Config) (*redis.Client, error) {
	_r := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Db,
	})

	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := _r.Ping(c).Result(); err != nil {
		return nil, err
	}
	return _r, nil
}
