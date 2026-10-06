package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
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
	texts := []string{"软件定义网络", "可编程数据平面", "向量检索"}
	vectors, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return err
	}
	if len(vectors) != len(texts) {
		return fmt.Errorf("向量数量与输入文本数量不一致")
	}
	for i, vector := range vectors {
		fmt.Printf("文本 %d 的向量维度：%d\n", i+1, len(vector))
	}
	return nil
}
