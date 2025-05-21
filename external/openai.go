package external

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIService struct {
	client openai.Client
	once   sync.Once
}

func (o *OpenAIService) GetOpenAIClient() openai.Client {
	o.once.Do(func() {
		openAIKey := os.Getenv("OPENAI_API_KEY")
		if openAIKey == "" {
			log.Fatalf("OPENAI_API_KEY is not set")
		}

		o.client = openai.NewClient(
			option.WithAPIKey(openAIKey),
		)
	})

	return o.client
}

func (o *OpenAIService) Ask(prompt string) string {
	chatCompletion, err := o.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		panic(err.Error())
	}

	return chatCompletion.Choices[0].Message.Content
}
