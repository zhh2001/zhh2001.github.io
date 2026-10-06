package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	ri "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/schema"
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
	// 按所选嵌入模型的实际输出维度创建新索引。
	probe, err := embedder.EmbedStrings(ctx, []string{"维度检查"})
	if err != nil {
		return err
	}
	if len(probe) != 1 || len(probe[0]) == 0 {
		return fmt.Errorf("嵌入模型未返回有效向量")
	}
	dimension := len(probe[0])
	err = client.FTCreate(ctx, "eino_notes", &goredis.FTCreateOptions{
		OnHash: true,
		Prefix: []any{"eino:note:"},
	},
		&goredis.FieldSchema{FieldName: "content", FieldType: goredis.SearchFieldTypeText},
		&goredis.FieldSchema{
			FieldName: "vector_content",
			FieldType: goredis.SearchFieldTypeVector,
			VectorArgs: &goredis.FTVectorArgs{
				FlatOptions: &goredis.FTFlatOptions{
					Type: "FLOAT32", Dim: dimension, DistanceMetric: "COSINE",
				},
			},
		},
	).Err()
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}
	indexer, err := ri.NewIndexer(ctx, &ri.IndexerConfig{
		Client:    client,
		KeyPrefix: "eino:note:",
		Embedding: embedder,
	})
	if err != nil {
		return err
	}
	ids, err := indexer.Store(ctx, []*schema.Document{
		{ID: "p4", Content: "P4 用于描述可编程数据平面的包处理逻辑。"},
		{ID: "sdn", Content: "SDN 将网络控制逻辑与数据转发功能分离。"},
		{ID: "int", Content: "INT 在数据包中携带网络遥测信息。"},
	})
	if err != nil {
		return err
	}
	fmt.Printf("索引维度：%d，文档 ID：%v\n", dimension, ids)
	return nil
}
