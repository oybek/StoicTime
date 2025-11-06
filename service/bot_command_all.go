package service

import (
	"context"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"oybek.io/sigma/model"
	"oybek.io/sigma/rdb"
)

func (b *Bot) handleCommandAll(tg *gotgbot.Bot, tgctx *ext.Context) error {
	ctx := context.Background()

	chat := tgctx.EffectiveChat

	b.typing(chat)

	acts, err := b.actStorage.FindAct(ctx, rdb.FindActArg{UserID: chat.Id})
	if err != nil {
		_, err = b.tg.SendMessage(chat.Id, "Что-то пошло не так, обратитесь @wolfodav", nil)
		return err
	}

	_, err = b.tg.SendMessage(chat.Id, actsToText(acts), nil)
	return err
}

func actsToText(acts []model.Act) string {
	if len(acts) == 0 {
		return "Нет занятий"
	}
	s := ""
	for i, act := range acts {
		s += fmt.Sprintf("%d. %s\n", i+1, act.Name)
	}
	return s
}
