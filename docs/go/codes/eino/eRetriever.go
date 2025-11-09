package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	rr "github.com/cloudwego/eino-ext/components/retriever/redis"
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

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	initRedis()

	ctx := context.Background()
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("EMBEDDER"),
	})
	if err != nil {
		panic(err)
	}

	retriever, err := rr.NewRetriever(ctx, &rr.RetrieverConfig{
		Client:       redisClient,
		Index:        "my_index",
		VectorField:  "vector_content",
		ReturnFields: []string{"content", "vector_content"},
		TopK:         2,
		Embedding:    embedder,
	})
	if err != nil {
		panic(err)
	}

	docs, err := retriever.Retrieve(ctx, "张恒华")
	if err != nil {
		panic(err)
	}

	for index, doc := range docs {
		fmt.Println(index, doc)
	}

}
