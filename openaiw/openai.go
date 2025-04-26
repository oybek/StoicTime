package openaiw

import "github.com/sashabaranov/go-openai"

type OpenAIWrapper struct {
	client *openai.Client
}

func NewOpenAIWrapper(client *openai.Client) *OpenAIWrapper {
	return &OpenAIWrapper{
		client: client,
	}
}
