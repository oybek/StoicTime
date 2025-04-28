package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/sashabaranov/go-openai"
	"oybek.io/kerege/config"
	"oybek.io/kerege/openaiw"
	"oybek.io/kerege/postgres"
	"oybek.io/kerege/service"
)

//go:embed assets/taxcode.txt
var taxCode string

func main() {
	ctx := context.Background()

	// load the configs
	theConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error create config: %v", err)
	}

	// di
	client := openai.NewClient(theConfig.OpenAIToken)
	openAIWrapper := openaiw.NewOpenAIWrapper(client)

	pgConn, err := pgx.Connect(ctx, theConfig.PGDSN)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer pgConn.Close(ctx)

	documentStorage := postgres.NewDocumentStorage(pgConn)

	embeddingService := service.NewEmbeddingService(documentStorage, openAIWrapper)
	if false {
		embeddingService.GetAndStoreEmbeddings(ctx, "hello")
	}

	bot, err := service.NewBot(
		theConfig.BotToken,
		client,
		documentStorage,
		openAIWrapper,
	)
	if err != nil {
		log.Fatalf("failed to create bot client: %v", err)
	}

	// launch
	go bot.Start()

	// listen for ctrl+c signal from terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping the bot...")
}
