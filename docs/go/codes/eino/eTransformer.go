package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"
)

var redisClient *goredis.Client

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: true, // 是否在输出的内容中移除标题行
	})
	if err != nil {
		panic(err)
	}

	docs := []*schema.Document{
		{
			ID: "doc1",
			Content: `# 文档标题

这是介绍部分的内容。

## 第一章

这是第一章的内容。

### 1.1 节

这是 1.1 节的内容。

## 第二章

这是第二章的内容。`,
		},
	}

	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result.Content)
		for k, v := range result.MetaData {
			fmt.Println("\t", k, v)
		}
	}
}
