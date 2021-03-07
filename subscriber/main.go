package main

import (
	"context"
	"line/health/service"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	service.LoadConfig()
}

func main() {
	// Uncomment below line incase local run
	//os.Setenv("RDB", "3")

	redisClient, err := service.ConnectRedis()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SLAVE READY")
	// assign worker
	workers := viper.GetInt("app.worker")
	var ctx = context.Background()
	wg := new(sync.WaitGroup)
	in := make(chan string, 2*workers)

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
				// http request
				result := service.FetchHTTP(msg.Payload)
				// set value to redis
				service.SetResult(ctx, redisClient, result)
			}
		}()
	}
	close(in)
	wg.Wait()
}
