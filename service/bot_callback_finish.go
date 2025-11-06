package service

import (
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (b *Bot) handleCallbackFinish(tg *gotgbot.Bot, tgCtx *ext.Context) error {
	tg.AnswerCallbackQuery(tgCtx.CallbackQuery.Id, nil)
	ctx := context.Background()
	chat := tgCtx.EffectiveChat
	err := b.finishLastActLog(ctx, chat.Id)
	return err
}
