package main

import (
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	ctx := context.Background()

	// 设置相关环境变量 COZELOOP_WORKSPACE_ID、COZELOOP_API_TOKEN
	client, err := cozeloop.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close(ctx)
	// 在服务 init 时 once 调用
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	// ...
}