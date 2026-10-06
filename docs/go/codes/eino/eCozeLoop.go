package main

import (
	"context"
	"fmt"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/coze-dev/cozeloop-go"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if os.Getenv("COZELOOP_WORKSPACE_ID") == "" || os.Getenv("COZELOOP_API_TOKEN") == "" {
		return fmt.Errorf("请设置 COZELOOP_WORKSPACE_ID 和 COZELOOP_API_TOKEN")
	}
	client, err := cozeloop.NewClient()
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		client.Close(closeCtx)
	}()
	// 进程初始化时注册一次，不要在每次请求时重复追加。
	callbacks.AppendGlobalHandlers(ccb.NewLoopHandler(client))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	chain := compose.NewChain[string, string]()
	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		return strings.ToUpper(input), nil
	}), compose.WithNodeName("upper"))
	runnable, err := chain.Compile(ctx)
	if err != nil {
		return err
	}
	output, err := runnable.Invoke(ctx, "eino")
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
