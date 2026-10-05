package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := rdb.Pipeline()
	incr := pipe.Incr(ctx, "pipeline-demo:counter")
	expire := pipe.Expire(ctx, "pipeline-demo:counter", time.Hour)
	if _, err := pipe.Exec(ctx); err != nil {
		panic(err)
	}
	// Exec 已经完成，才读取各命令结果
	value, err := incr.Result()
	if err != nil {
		panic(err)
	}
	if err := expire.Err(); err != nil {
		panic(err)
	}
	fmt.Println(value)
}
