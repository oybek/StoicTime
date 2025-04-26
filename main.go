package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
	"oybek.io/kerege/config"
	"oybek.io/kerege/openaiw"
)

const text = `
Кереге это умный консультант по кодексам и законам республики Кыргызстан
`

func main() {
	ctx := context.Background()

	theConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error create config: %v", err)
	}

	var openAIWrapper *openaiw.OpenAIWrapper
	{
		client := openai.NewClient(theConfig.OpenAIToken)
		openAIWrapper = openaiw.NewOpenAIWrapper(client)
	}

	embedding, _ := openAIWrapper.GetEmbedding(ctx, text)

	fmt.Printf("%#v\n", embedding)
}
