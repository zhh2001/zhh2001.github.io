package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"log"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	noteTool := CreateTool()
	result, err := noteTool.InvokableRun(context.Background(), `{"name":"P4"}`)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=技术名称"`
}

func GetNote(ctx context.Context, params *InputParams) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if params == nil || strings.TrimSpace(params.Name) == "" {
		return "", fmt.Errorf("技术名称不能为空")
	}
	notes := map[string]string{
		"p4":      "https://zhh2001.github.io/sdn/p4",
		"int":     "https://zhh2001.github.io/sdn/int",
		"mininet": "https://zhh2001.github.io/sdn/mininet",
		"iperf":   "https://zhh2001.github.io/sdn/iperf",
	}
	url, ok := notes[strings.ToLower(strings.TrimSpace(params.Name))]
	if !ok {
		return "", fmt.Errorf("未找到对应笔记")
	}
	return url, nil
}

func CreateTool() tool.InvokableTool {
	return utils.NewTool(&schema.ToolInfo{
		Name: "get_note",
		Desc: "根据技术名称获取学习笔记链接",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"name": {Type: schema.String, Desc: "技术名称", Required: true},
		}),
	}, GetNote)
}
