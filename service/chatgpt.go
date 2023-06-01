package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/JabinGP/demo-chatroom/config"
	"github.com/JabinGP/demo-chatroom/infra/logger"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"strings"
)

// MessageService message service
type ChatGptService struct {
	systemUser string
	logger     *logger.CustZeroLogger
}

var client *openai.Client

func (c *ChatGptService) Init() {
	apiBase := config.Viper.GetString("openai.apiBase")
	apiKey := config.Viper.GetString("openai.apiKey")
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = apiBase
	c.logger.Info().Msgf("query chatgpt from azure %v, apiBase %s apiKey %s", strings.Contains(apiBase, "azure.com"), apiBase, apiKey)
	if strings.Contains(apiBase, "azure.com") {
		c.logger.Info().Msgf("query chatgpt apiBase %s apiKey %s", apiBase, apiKey)
		//config := openai.DefaultAzureConfig(apiKey, apiBase)
		config.APIVersion = "2023-05-15" // optional update to latest API version
		config.APIType = openai.APITypeAzure
		//If you use a deployment name different from the model name, you can customize the AzureModelMapperFunc function
		config.AzureModelMapperFunc = func(model string) string {
			azureModelMapping := map[string]string{
				"gpt-3.5-turbo": "gpt-35-turbo", // "your gpt-3.5-turbo deployment name",
			}
			return azureModelMapping[model]
		}
	}

	client = openai.NewClientWithConfig(config)
}

func (c *ChatGptService) Embeddings(input string) {
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

func (c *ChatGptService) Ask(question string, isStream bool, out io.Writer) error {

	//req := openai.ChatCompletionRequest{
	//	Model: conf.Model,
	//	Messages: []openai.ChatCompletionMessage{
	//		{Role: openai.ChatMessageRoleSystem, Content: c.globalConf.LookupPrompt(conf.Prompt)},
	//		{Role: openai.ChatMessageRoleUser, Content: question},
	//	},
	//	MaxTokens:   conf.MaxTokens,
	//	Temperature: conf.Temperature,
	//	N:           1,
	//}

	c.logger.Info().Msgf("question %s, stream %v", question, isStream)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: c.systemUser},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	}
	if isStream {
		req.Stream = true
		stream, err := client.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			c.logger.Error().Msgf("query chatgpt error %v", err)
			return err
		}
		defer stream.Close()
		for {
			resp, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					_, _ = fmt.Fprintln(out)

					break
				}
				return err
			}
			content := resp.Choices[0].Delta.Content
			c.logger.Info().Msgf("stream content %v", resp)
			_, _ = fmt.Fprint(out, content)
		}
	} else {
		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			return err
		}
		content := resp.Choices[0].Message.Content
		_, _ = fmt.Fprintln(out, content)
	}
	return nil
}
