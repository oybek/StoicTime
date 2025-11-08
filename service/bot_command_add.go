package service

import (
	"context"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"oybek.io/sigma/model"
)

func refineActName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (b *Bot) typing(chat *gotgbot.Chat) {
	b.tg.SendChatAction(chat.Id, "typing", nil)
}

func (b *Bot) handleCommandAdd(tg *gotgbot.Bot, tgctx *ext.Context) error {
	ctx := context.Background()

	chat := tgctx.EffectiveChat
	text := tgctx.EffectiveMessage.Text

	b.typing(chat)

	actName := refineActName(strings.TrimPrefix(text, "/add"))
	if actName == "" {
		return b.handleCommandHelp(tg, tgctx)
	}

	err := b.actStorage.CreateAct(ctx, model.Act{UserID: chat.Id, Name: actName})
	if err != nil {
		_, err = b.tg.SendMessage(chat.Id, "Такое занятие уже есть", nil)
		return err
	}

	_, err = b.tg.SendMessage(chat.Id, "Новое занятие: "+actName, nil)
	return err
}
