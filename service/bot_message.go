package service

import (
	"context"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"oybek.io/sigma/model"
	"oybek.io/sigma/rdb"
	"oybek.io/sigma/texts"
)

func (b *Bot) onMessage(tg *gotgbot.Bot, tgCtx *ext.Context) error {
	ctx := context.Background()

	chat := tgCtx.EffectiveChat
	text := tgCtx.EffectiveMessage.Text

	b.typing(chat)

	actName, err := b.checkActExists(ctx, chat.Id, text)
	if err != nil || actName == "" {
		return err
	}

	err = b.finishLastActLog(ctx, chat.Id)
	if err != nil {
		return err
	}

	return b.createActLog(ctx, chat.Id, actName)
}

func (b *Bot) checkActExists(ctx context.Context, chatID int64, text string) (string, error) {
	actName := strings.TrimSpace(text)
	acts, err := b.actStorage.FindAct(ctx, rdb.FindActArg{
		UserID: chatID,
		Name:   actName,
	})
	if err != nil {
		return "", err
	}
	if len(acts) == 0 {
		_, err := b.tg.SendMessage(chatID, texts.ActMissing, nil)
		return "", err
	}
	return actName, nil
}

func (b *Bot) finishLastActLog(ctx context.Context, chatID int64) error {
	actLog, err := b.actLogStorage.FindActLog(ctx, rdb.FindActLogArg{UserID: chatID, Active: true})
	if err != nil || len(actLog) == 0 {
		return err
	}

	actLog[0].EndTime = time.Now().UTC()
	err = b.actLogStorage.UpdateActLog(ctx, actLog[0])
	if err != nil {
		return err
	}

	b.tg.EditMessageText(actLog[0].Text(), &gotgbot.EditMessageTextOpts{
		ChatId:    chatID,
		MessageId: actLog[0].MessageID,
	})
	return nil
}

func (b *Bot) createActLog(ctx context.Context, chatID int64, actName string) error {
	actLog := model.ActLog{
		UserID:    chatID,
		Name:      actName,
		StartTime: time.Now().UTC(),
	}
	message, err := b.tg.SendMessage(chatID, actLog.Text(), &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				gotgbot.InlineKeyboardButton{
					Text:         "Завершить",
					CallbackData: "/finish",
				},
			}},
		},
	})
	if err != nil {
		return err
	}

	actLog.MessageID = message.MessageId
	err = b.actLogStorage.CreateActLog(ctx, actLog)
	if err != nil {
		b.tg.DeleteMessage(chatID, message.MessageId, nil)
		return err
	}

	return nil
}
