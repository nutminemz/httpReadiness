package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	"line/health/service"

	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	service.LoadConfig()
}

func main() {
	// Uncomment 2 below line in case local run
	//os.Setenv("SLAVE1", "3")
	//os.Setenv("SLAVE2", "3")

	redisClient, err := service.ConnectRedis()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("START PROCESS")

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

	// line counter
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
