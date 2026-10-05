package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"127.0.0.1:7101", "127.0.0.1:7102", "127.0.0.1:7103"},
	})
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Set(ctx, "cluster-demo:name", "zhh", time.Minute).Err(); err != nil {
		panic(err)
	}
	value, err := rdb.Get(ctx, "cluster-demo:name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

	// 回调可能并发执行，覆盖客户端当前已知的主节点和副本
	err = rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}
}
