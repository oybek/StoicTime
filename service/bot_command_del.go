package service

import (
	"context"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"oybek.io/sigma/rdb"
)

func (b *Bot) handleCommandDel(tg *gotgbot.Bot, tgctx *ext.Context) error {
	ctx := context.Background()

	chat := tgctx.EffectiveChat
	text := tgctx.EffectiveMessage.Text

	b.typing(chat)

	actName := refineActName(strings.TrimPrefix(text, "/del"))
	if actName == "" {
		return b.handleCommandHelp(tg, tgctx)
	}

	err := b.actStorage.DeleteAct(ctx, rdb.DeleteActArg{UserID: chat.Id, Name: actName})
	if err != nil {
		_, err = b.tg.SendMessage(chat.Id, "Такого занятия нет", nil)
		return err
	}

	_, err = b.tg.SendMessage(chat.Id, "Удалено занятие: "+actName, nil)
	return err
}
