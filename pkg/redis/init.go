package redis

import (
	"context"
	"log"
	"time"
	"tj/config"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Connect() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: config.Cfg.RedisAddr,
		DB:   0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := Rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	log.Println("Redis connected")
}
