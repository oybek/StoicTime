package service

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var helpText = strings.TrimSpace(`
/add <занятие> - добавить занятие
/del <занятие> - удалить занятие
/all - список занятий
<занятие> - начать занятие
`)

func (b *Bot) handleCommandHelp(_ *gotgbot.Bot, tgctx *ext.Context) error {
	chat := tgctx.EffectiveChat
	_, err := b.tg.SendMessage(chat.Id, helpText, nil)
	return err
}
