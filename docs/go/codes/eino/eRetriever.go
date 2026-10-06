package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	rr "github.com/cloudwego/eino-ext/components/retriever/redis"
	goredis "github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apiKey, modelID := os.Getenv("ARK_API_KEY"), os.Getenv("EMBEDDER")
	if apiKey == "" || modelID == "" {
		return fmt.Errorf("请设置 ARK_API_KEY 和 EMBEDDER")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	apiType := ark.APITypeText
	if value := os.Getenv("EMBEDDER_API_TYPE"); value != "" {
		apiType = ark.APIType(value)
		if apiType != ark.APITypeText && apiType != ark.APITypeMultiModal {
			return fmt.Errorf("EMBEDDER_API_TYPE 应为 text_api 或 multi_modal_api")
		}
	}
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  apiKey,
		Model:   modelID,
		BaseURL: os.Getenv("ARK_BASE_URL"),
		APIType: &apiType,
	})
	if err != nil {
		return err
	}
	address := os.Getenv("REDIS_ADDR")
	if address == "" {
		return fmt.Errorf("请设置 REDIS_ADDR")
	}
	client := goredis.NewClient(&goredis.Options{
		Addr:     address,
		Password: os.Getenv("REDIS_PASSWORD"),
		Protocol: 2,
	})
	defer client.Close()
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}
	retriever, err := rr.NewRetriever(ctx, &rr.RetrieverConfig{
		Client:       client,
		Index:        "eino_notes",
		VectorField:  "vector_content",
		ReturnFields: []string{"content"},
		TopK:         2,
		Embedding:    embedder,
	})
	if err != nil {
		return err
	}
	docs, err := retriever.Retrieve(ctx, "什么是可编程数据平面？")
	if err != nil {
		return err
	}
	for _, doc := range docs {
		fmt.Printf("%s: %s\n", doc.ID, doc.Content)
	}
	return nil
}
