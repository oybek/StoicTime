package service

import (
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/sashabaranov/go-openai"
)

type Bot struct {
	bot                 *gotgbot.Bot
	client              *openai.Client
	documentStorage     DocumentStorage
	embeddingCalculator EmbeddingCalculator
}

func NewBot(
	botToken string,
	client *openai.Client,
	documentStorage DocumentStorage,
	embeddingCalculator EmbeddingCalculator,
) (*Bot, error) {
	bot, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:                 bot,
		client:              client,
		documentStorage:     documentStorage,
		embeddingCalculator: embeddingCalculator,
	}, nil
}

func (b *Bot) Start() error {
	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewMessage(message.Text, b.onMessage))

	// Start receiving updates.
	err := updater.StartPolling(b.bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 15,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		return err
	}

	log.Printf("%s has been started...\n", b.bot.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
	return nil
}
