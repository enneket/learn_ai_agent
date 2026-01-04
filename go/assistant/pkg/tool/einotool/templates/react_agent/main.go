package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Printf("读取.env文件失败：%v\n", err)
		return
	}

	ctx := context.Background()
	reactAgent, err := NewAgent(ctx)
	if err != nil {
		panic(err)
	}

	arg := flag.Arg(0)
	if arg == "" {
		panic("message is required, eg: ./llm -model=ep-xxxx -apikey=xxx 'do you know cloudwego?'")
	}

	sr, err := reactAgent.Stream(ctx, []*schema.Message{
		schema.UserMessage(arg),
	}, agent.WithComposeOptions(compose.WithCallbacks(LogCallback())))
	if err != nil {
		panic(err)
	}

	for {
		msg, err := sr.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Print(msg.Content)
	}
	fmt.Printf("\n\n=== %sFINISHED%s ===\n\n", green, reset)
}

func NewAgent(ctx context.Context) (*react.Agent, error) {

	// 初始化模型
	model, err := PrepareModel(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化各种 tool
	tools, err := PrepareTools(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化 agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model: model,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func PrepareModel(ctx context.Context) (model.ChatModel, error) {
	qwenModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		Model:   os.Getenv("QWEEN_MODEL_NAME"), // replace with your model
		APIKey:  os.Getenv("QWEEN_API_KEY"),    // replace with your api key
		BaseURL: os.Getenv("QWEEN_BASE_URL"),
	})
	if err != nil {
		return nil, err
	}
	return qwenModel, nil
}

func PrepareTools(ctx context.Context) ([]tool.BaseTool, error) {
	ddg, err := duckduckgo.NewTextSearchTool(ctx, &duckduckgo.Config{})
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{ddg}, nil
}

// log with color
var (
	green = "\033[32m"
	reset = "\033[0m"
)

func LogCallback() callbacks.Handler {
	builder := callbacks.NewHandlerBuilder()
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Printf("%s[view]%s: start [%s:%s:%s]\n", green, reset, info.Component, info.Type, info.Name)
		return ctx
	})
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Printf("%s[view]%s: end [%s:%s:%s]\n", green, reset, info.Component, info.Type, info.Name)
		return ctx
	})
	return builder.Build()
}
