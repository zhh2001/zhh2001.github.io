package main

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type Note struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=技术名"`
}

func GetNote(ctx context.Context, params *InputParams) (string, error) {
	NoteSet := []Note{
		{Name: "P4", Url: "https://zhh2001.github.io/sdn/p4"},
		{Name: "INT", Url: "https://zhh2001.github.io/sdn/int"},
		{Name: "Mininet", Url: "https://zhh2001.github.io/sdn/mininet"},
		{Name: "iPerf", Url: "https://zhh2001.github.io/sdn/iperf"},
	}
	for _, note := range NoteSet {
		if strings.ToLower(params.Name) == strings.ToLower(note.Name) {
			return note.Url, nil
		}
	}
	return "", nil
}

func CreateTool() tool.InvokableTool {
	GetNoteTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_note",
		Desc: "获取笔记的链接",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": {
					Type:     schema.String,
					Desc:     "技术名",
					Required: true,
				},
			},
		),
	}, GetNote)
	return GetNoteTool
}
