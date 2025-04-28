package service

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const el = "\n"

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendMessage(
		ctx.EffectiveChat.Id,
		fmt.Sprintf(
			"Здравствуйте %s!"+el+
				"Я знаю налоговый кодекс Кыргызстана наизусть"+el+
				"Задавайте мне любые вопросы", ctx.EffectiveUser.FirstName,
		),
		&gotgbot.SendMessageOpts{},
	)
	if err != nil {
		return err
	}

	return nil
}
