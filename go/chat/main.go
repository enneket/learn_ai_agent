package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("读取.env文件失败：%v\n", err)
		return
	}

	ctx := context.Background()

	// 使用模版创建messages
	log.Printf("===create messages===\n")
	messages := createMessagesFromTemplate()
	log.Printf("messages: %+v\n\n", messages)

	// 创建llm
	log.Printf("===create llm===\n")
	cm := createQwenChatModel(ctx)
	log.Printf("create llm success\n\n")

	log.Printf("===llm generate===\n")
	result := generate(ctx, cm, messages)
	log.Printf("result: %+v\n\n", result)

	log.Printf("===llm stream generate===\n")
	streamResult := stream(ctx, cm, messages)
	reportStream(streamResult)
}

func createQwenChatModel(ctx context.Context) model.ToolCallingChatModel {
	key := os.Getenv("QWEEN_API_KEY")
	modelName := os.Getenv("QWEEN_MODEL_NAME")
	baseUrl := os.Getenv("QWEEN_BASE_URL")
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		APIKey:  key,
		Model:   modelName,
		BaseURL: baseUrl,
	})
	if err != nil {
		panic(err)
	}
	return chatModel
}

func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
}

func createMessagesFromTemplate() []*schema.Message {
	template := createTemplate()

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	return messages
}

func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()

	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF { // 流式输出结束
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		log.Printf("message[%d]: %+v\n", i, message)
		i++
	}
}

func generate(ctx context.Context, llm model.ToolCallingChatModel, in []*schema.Message) *schema.Message {
	result, err := llm.Generate(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

func stream(ctx context.Context, llm model.ToolCallingChatModel, in []*schema.Message) *schema.StreamReader[*schema.Message] {
	result, err := llm.Stream(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}
