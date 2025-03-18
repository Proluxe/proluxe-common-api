package services

import (
	"context"
	"log"

	"github.com/sashabaranov/go-openai"
)

type OpenAIConfig struct {
	OpenAIToken string
}

type LLM struct {
	Config  OpenAIConfig
	Client  *openai.Client
	Context context.Context
}

func InitOpenAI(apiKey string) LLM {
	config := &OpenAIConfig{
		OpenAIToken: apiKey,
	}

	ai := LLM{
		Config:  *config,
		Client:  openai.NewClient(config.OpenAIToken),
		Context: context.Background(),
	}

	return ai
}

func (ai *LLM) Ask(prompt string) (string, error) {
	resp, err := ai.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, err
}
