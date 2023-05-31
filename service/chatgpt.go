package service

import (
	"context"
	"fmt"
	"github.com/JabinGP/demo-chatroom/config"
	openai "github.com/sashabaranov/go-openai"
)

// MessageService message service
type ChatGptService struct {
}

var client *openai.Client

func init() {
	config := openai.DefaultAzureConfig(config.Viper.GetString("openai.apiKey"), config.Viper.GetString("openai.apiBase"))
	config.APIVersion = "2023-05-15" // optional update to latest API version

	//If you use a deployment name different from the model name, you can customize the AzureModelMapperFunc function
	config.AzureModelMapperFunc = func(model string) string {
		azureModelMapping := map[string]string{
			"gpt-3.5-turbo": "gpt-35-turbo", // "your gpt-3.5-turbo deployment name",
		}
		return azureModelMapping[model]
	}

	client = openai.NewClientWithConfig(config)
}

func (*ChatGptService) Embeddings(input string) {
	resp, err := client.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Input: []string{input},
			Model: openai.AdaEmbeddingV2,
		})

	if err != nil {
		fmt.Printf("CreateEmbeddings error: %v\n", err)
		return
	}

	vectors := resp.Data[0].Embedding // []float32 with 1536 dimensions

	fmt.Println(vectors[:10], "...", vectors[len(vectors)-10:])
}

func (*ChatGptService) ChatGPT(input string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	fmt.Println(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content, nil
}
