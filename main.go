package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sashabaranov/go-openai"
	"oybek.io/kerege/config"
	"oybek.io/kerege/openaiw"
)

const text = `
Кереге это умный консультант по кодексам и законам республики Кыргызстан
`

func main() {
	ctx := context.Background()

	//
	theConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error create config: %v", err)
	}

	//
	var openAIWrapper *openaiw.OpenAIWrapper
	{
		client := openai.NewClient(theConfig.OpenAIToken)
		openAIWrapper = openaiw.NewOpenAIWrapper(client)
	}

	pgConn, err := pgx.Connect(ctx, theConfig.PGDSN)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer pgConn.Close(ctx)

	embedding, _ := openAIWrapper.GetEmbedding(ctx, text)

	fmt.Printf("%#v\n", embedding)
}
