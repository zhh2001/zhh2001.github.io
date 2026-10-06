package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers:     map[string]string{"#": "h1", "##": "h2", "###": "h3"},
		TrimHeaders: true,
	})
	if err != nil {
		return err
	}
	docs := []*schema.Document{{
		ID:      "network-note",
		Content: "# 网络笔记\n\n介绍网络体系结构。\n\n## SDN\n\n控制平面与数据平面分离。\n\n### P4\n\n描述数据平面的包处理逻辑。\n\n## INT\n\n记录转发路径上的遥测信息。",
	}}
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result.Content)
		fmt.Printf("标题：%v / %v / %v\n", result.MetaData["h1"], result.MetaData["h2"], result.MetaData["h3"])
	}
	return nil
}
