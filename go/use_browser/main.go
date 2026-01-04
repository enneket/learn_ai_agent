package main

import (
	"context"
	"fmt"
	"learn_ai_agent/utils/gptr"
	"log"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
)

func main() {
	ctx := context.Background()

	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		log.Fatal(err)
	}

	info, err := but.Info(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(info)

	url := "https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/tool/tool_browseruse/"
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	result, err = but.Execute(&browseruse.Param{
		Action: browseruse.ActionExtractContent,
		Goal:   gptr.Of("BrowserUse Tool"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	time.Sleep(10 * time.Second)
	but.Cleanup()
}
