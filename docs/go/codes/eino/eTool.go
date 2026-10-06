package main

import (
	"context"
	"encoding/json"
	"fmt"
	req "github.com/cloudwego/eino-ext/components/tool/httprequest/get"
	"log"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	getTool, err := req.NewTool(ctx, &req.Config{
		Headers:    map[string]string{"User-Agent": "EinoNotesExample"},
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	})
	if err != nil {
		return err
	}
	arguments, err := json.Marshal(&req.GetRequest{URL: "https://zhh2001.github.io/sitemap.xml"})
	if err != nil {
		return err
	}
	result, err := getTool.InvokableRun(ctx, string(arguments))
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
