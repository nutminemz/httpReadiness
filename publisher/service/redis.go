package service

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// ConnectRedis :: connect to redis
func ConnectRedis() (c *redis.Client, e error) {
	log.Println("Waiting.. for redis")
	// waiting for redis startup
	time.Sleep(3 * time.Second)
	// Create a new Redis Client
	redisAddress := viper.GetString("redis.host")
	redisPassword := viper.GetString("redis.password")
	redisRetry := viper.GetInt("redis.retry")
	redisClient := redis.NewClient(&redis.Options{
		Addr:       redisAddress,
		Password:   redisPassword,
		MaxRetries: redisRetry,
	})

	// Ping the Redis server and check if any errors occured
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 7 seconds and wait for Redis to initialize
		log.Println("Waiting.. for redis")
		time.Sleep(7 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
	return redisClient, err
}
