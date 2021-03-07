package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
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
	log.Println("Waiting.. for redis")
	// Uncomment 2 below line incase local run
	//os.Setenv("SLAVE1", "3")
	//os.Setenv("SLAVE2", "3")

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
	// Start time Rec
	start := time.Now()

	// Read file
	inputPath := viper.GetString("input.path")
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	parser := csv.NewReader(file)

	// Generate a new background context that  we will use
	ctx := context.Background()
	// init result
	redisClient.Set(ctx, "success", 0, 0)
	redisClient.Set(ctx, "fail", 0, 0)
	i := 0
	// Loop till end of file
	for {
		record, err := parser.Read()
		if err == io.EOF {
			if err != nil {
				log.Println(err)
			}
			break
		}
		if err != nil {
			log.Println(err)
		}
		// Load balancer between 2 subscribers
		if i%2 == 0 {
			err = redisClient.Publish(ctx, os.Getenv("SLAVE1"), record[0]).Err()
		} else {
			err = redisClient.Publish(ctx, os.Getenv("SLAVE2"), record[0]).Err()
		}

		if err != nil {
			log.Println(err)
		}
		i++
	}

	// print summary result evert 1 sec
	for {
		<-time.After(1 * time.Second)
		s, _ := redisClient.Get(ctx, "success").Int()
		f, _ := redisClient.Get(ctx, "fail").Int()
		p := s + f
		log.Println("total(", p, "/", i, ") success:", s, " fail:", f)
		if p == i {
			// Stop timmer, print result, call LINE API
			elapsed := time.Since(start)
			log.Println("EMPTY QUEUE")
			log.Printf("Process took %s ", elapsed)
			log.Println("POST to LINE API")
			service.LineResult(i, s, f, int64(elapsed))
			log.Println("END PROCESS")
			break
		}
	}
}
