package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// ConnectRedis :: connect to redis
func ConnectRedis() (c *redis.Client, e error) {
	// Create a new Redis Client
	redisAddress := viper.GetString("redis.host")
	redisPassword := viper.GetString("redis.password")
	redisRetry := viper.GetInt("redis.retry")
	poolSize := viper.GetInt("redis.poolsize")
	minIdieConns := viper.GetInt("redis.minidieconns")
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisAddress,
		Password:     redisPassword,
		MaxRetries:   redisRetry,
		PoolSize:     poolSize,
		MinIdleConns: minIdieConns,
	})

	// Ping the Redis server and check if any errors occured
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
	return redisClient, err
}
