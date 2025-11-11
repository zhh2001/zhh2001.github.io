package main

import (
	"context"
	"hash/crc32"

	"github.com/golang/groupcache/consistenthash"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			// shardName => host:port
			"shard1": "localhost:7000",
			"shard2": "localhost:7001",
			"shard3": "localhost:7002",
		},

		// 手动设置连接节点
		NewClient: func(opt *redis.Options) *redis.Client {
			opt.Username = "dev"
			opt.Password = "pwd"

			return redis.NewClient(opt)
		},

		// 更改为其他 Hash 算法，go-redis 默认使用 Rendezvous Hash 算法将 Key 分布到多个节点上
		NewConsistentHash: func(shards []string) redis.ConsistentHash {
			return consistenthash.New(100, crc32.ChecksumIEEE)
		},
	})

	ctx := context.Background()

	// 执行命令
	if err := rdb.Set(ctx, "foo", "bar", 0).Err(); err != nil {
		panic(err)
	}

	// 遍历每个节点
	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}
}
