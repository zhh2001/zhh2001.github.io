package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 初始化嵌入器
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("EMBEDDER"),
	})
	if err != nil {
		panic(err)
	}

	input := []string{
		"你好，泥嚎",
		"微服务",
		"大模型",
	}
	embeddings, err := embedder.EmbedStrings(ctx, input)
	if err != nil {
		panic(err)
	}
	for i, embedding := range embeddings {
		fmt.Printf("文本 (%d) 的向量维度是 %d\n", i+1, len(embedding))
	}
}
