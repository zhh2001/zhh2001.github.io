package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	ri "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"
)

var redisClient *goredis.Client

func initRedis() {
	redisClient = goredis.NewClient(&goredis.Options{
		Addr:          os.Getenv("REDIS_ADDR"),
		Protocol:      2,
		UnstableResp3: true,
	})
}

func createIndex() {
	ctx := context.Background()

	keyPrefix := "eino:"    // keyPrefix should be the prefix of keys you write to redis and want to retrieve.
	indexName := "my_index" // indexName should be used in redis retriever.

	// schemas should match DocumentToHashes configured in IndexerConfig.
	schemas := []*goredis.FieldSchema{
		{
			FieldName: "content",
			FieldType: goredis.SearchFieldTypeText,
			Weight:    1,
		},
		{
			FieldName: "vector_content",
			FieldType: goredis.SearchFieldTypeVector,
			VectorArgs: &goredis.FTVectorArgs{
				FlatOptions: &goredis.FTFlatOptions{
					Type:           "FLOAT32", // BFLOAT16 / FLOAT16 / FLOAT32 / FLOAT64. BFLOAT16 and FLOAT16
					Dim:            2560,      // keeps same with dimensions of Embedding
					DistanceMetric: "COSINE",  // L2 / IP / COSINE
				},
				HNSWOptions: nil,
			},
		},
		{
			FieldName: "extra_field_number",
			FieldType: goredis.SearchFieldTypeNumeric,
		},
	}

	options := &goredis.FTCreateOptions{
		OnHash: true,
		Prefix: []any{keyPrefix},
	}

	result, err := redisClient.FTCreate(ctx, indexName, options, schemas...).Result()
	if err != nil && err.Error() != "Index already exists" {
		panic(err)
	}

	fmt.Println(result)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	initRedis()
	createIndex()

	ctx := context.Background()

	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("EMBEDDER"),
	})
	if err != nil {
		panic(err)
	}

	indexer, err := ri.NewIndexer(ctx, &ri.IndexerConfig{
		Client:    redisClient,
		KeyPrefix: "eino:",
		Embedding: embedder,
	})
	if err != nil {
		panic(err)
	}
	ids, err := indexer.Store(ctx, []*schema.Document{
		{
			ID:      "1",
			Content: "张恒华同学是上海工程技术大学的学生，研究方向是软件定义网络。",
			MetaData: map[string]any{
				"author": "Howard",
			},
		},
		{
			ID:      "2",
			Content: "你说得对，但是 PHP 是世界上最好的语言。",
			MetaData: map[string]any{
				"author": "Howard",
			},
		},
		{
			ID:      "3",
			Content: "张恒华的英文名为 Howard Zhang。",
			MetaData: map[string]any{
				"author": "Howard",
			},
		},
		{
			ID:      "4",
			Content: "张恒华同学今年 18 岁了。",
			MetaData: map[string]any{
				"author": "Howard",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(ids)
}
