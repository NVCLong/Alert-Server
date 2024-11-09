package redis

import (
	"crypto/tls"

	env "github.com/NVCLong/Alert-Server/bootstrap"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection() *redis.Client {
	redisOption := &redis.Options{
		Addr:      env.GetEnv(env.EnvRedisHost),
		Password:  env.GetEnv(env.EnvRedisAccessKey),
		TLSConfig: &tls.Config{},
	}
	client := redis.NewClient(redisOption)
	return client
}
