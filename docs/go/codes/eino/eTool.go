package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	req "github.com/cloudwego/eino-ext/components/tool/httprequest/get"
)

func main() {
	ctx := context.Background()

	config := &req.Config{
		Headers: map[string]string{
			"User-Agent": "MyCustomAgent",
		},
		HttpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: &http.Transport{},
		},
	}

	tool, err := req.NewTool(ctx, config)
	if err != nil {
		panic(err)
	}

	request := &req.GetRequest{
		URL: "https://zhh2001.github.io/sitemap",
	}

	jsonReq, err := sonic.Marshal(request)
	if err != nil {
		panic(err)
	}

	resp, err := tool.InvokableRun(ctx, string(jsonReq))
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
