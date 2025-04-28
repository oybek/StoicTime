package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/sashabaranov/go-openai"
)

func (b *Bot) onMessage(_ *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	text := ctx.EffectiveMessage.Text

	b.bot.SendChatAction(chat.Id, "typing", nil)

	aText, err := b.answerTheQuestion(context.Background(), text)
	if err != nil {
		return err
	}

	_, err = b.bot.SendMessage(ctx.EffectiveChat.Id, aText, nil)
	return err
}

func (b *Bot) answerTheQuestion(ctx context.Context, qText string) (string, error) {
	qEmbedding, err := b.embeddingCalculator.GetEmbedding(ctx, qText)
	if err != nil {
		return "", err
	}

	texts, err := b.documentStorage.Search(ctx, qEmbedding)
	if err != nil {
		return "", err
	}

	contextText := strings.Join(texts, "\n---\n")
	prompt := fmt.Sprintf("Контекст:\n%s\n\nВопрос: %s\n", contextText, qText)

	// Получаем ответ от GPT
	resp, err := b.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "Ты помощник, отвечающий по контексту. Все в рамках страны Кыргызстан. " +
				"Если в контексте нет нужной информации - советуй обратиться в юридическую фирму 'Кереге'"},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
