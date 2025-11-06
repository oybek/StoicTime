package service

import (
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type Bot struct {
	tg            *gotgbot.Bot
	clock         *Clock
	actStorage    ActStorage
	actLogStorage ActLogStorage
}

func NewBot(
	botToken string,
	clock *Clock,
	actStorage ActStorage,
	actLogStorage ActLogStorage,
) (*Bot, error) {
	bot, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		return nil, err
	}

	return &Bot{
		tg:            bot,
		clock:         clock,
		actStorage:    actStorage,
		actLogStorage: actLogStorage,
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
	dispatcher.AddHandler(handlers.NewCommand("start", b.handleCommandHelp))
	dispatcher.AddHandler(handlers.NewCommand("help", b.handleCommandHelp))
	dispatcher.AddHandler(handlers.NewCommand("add", b.handleCommandAdd))
	dispatcher.AddHandler(handlers.NewCommand("del", b.handleCommandDel))
	dispatcher.AddHandler(handlers.NewCommand("all", b.handleCommandAll))
	dispatcher.AddHandler(handlers.NewCommand("rep", b.handleCommandRep))
	dispatcher.AddHandler(handlers.NewMessage(message.Text, b.onMessage))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("/finish"), b.handleCallbackFinish))

	// Start receiving updates.
	err := updater.StartPolling(b.tg, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 15,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 20,
			},
		},
	})
	if err != nil {
		return err
	}

	log.Printf("%s has been started...\n", b.tg.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
	return nil
}
