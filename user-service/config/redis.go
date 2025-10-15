package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func (cfg Config) NewRedisClient() *redis.Client {
	// connect := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Memastikan Redis terhubung
	_, err := client.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}

	return client
}
