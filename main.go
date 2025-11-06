package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"oybek.io/sigma/config"
	"oybek.io/sigma/rdb"
	"oybek.io/sigma/service"
)

func main() {
	ctx := context.Background()

	theConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	rdb.Migrate(theConfig.PGURL)

	theRdb, err := rdb.NewRdb(ctx, theConfig.PGURL)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	bot, err := service.NewBot(
		theConfig.BotToken,
		theRdb,
		theRdb,
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
