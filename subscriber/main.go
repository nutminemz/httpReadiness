package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"line/health/service"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	// setup config
	viper.AddConfigPath("./conf")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.MergeConfig(strings.NewReader(viper.GetString("configs"))); err != nil {
		log.Panic(err.Error())
	} else {
		log.Println("loaded config " + viper.GetString("app.name"))
	}
	log.Println(viper.AllSettings())
}

func main() {
	// Uncomment below line incase local run
	//os.Setenv("RDB", "3")

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
	workers := viper.GetInt("app.worker")
	var ctx = context.Background()
	wg := new(sync.WaitGroup)
	in := make(chan string, 2*workers)

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
	log.Println("SLAVE READY")
	// Subscribe to the Topic given
	topic := redisClient.Subscribe(ctx, os.Getenv("RDB"))
	// Get the Channel to use
	channelSize := viper.GetInt("redis.channelsize")
	channel := topic.ChannelSize(channelSize)
	// Itterate any messages sent on the Topic
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range channel {
				service.SetResult(ctx, redisClient, service.FetchHTTP(msg))
			}
		}()
	}
	close(in)
	wg.Wait()
}
