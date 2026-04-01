package llm

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/spf13/viper"
)
import "github.com/openai/openai-go/v3"

type LLMClient struct {
	ApiKey   string
	Client   *openai.Client
	Model    string
	EndPoint string
}

func NewLLMClient(vc *viper.Viper) *LLMClient {
	apiKey := vc.GetString("llm.apikey")
	endpoint := vc.GetString("llm.endpoint")
	model := vc.GetString("llm.model")
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(endpoint),
	)
	return &LLMClient{
		ApiKey:   apiKey,
		Client:   &client,
		Model:    model,
		EndPoint: endpoint,
	}
}

func (lc *LLMClient) GetResponse(messages *Messages) (*openai.ChatCompletion, error) {
	chatCompletion, err := lc.Client.Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: messages.ToMessageParamUnion(),
			Model:    lc.Model,
		},
		option.WithJSONSet("stream", false),
		option.WithJSONSet("enable_thinking", false),
	)
	if err != nil {
		return nil, err
	}
	return chatCompletion, nil
}

func (lc *LLMClient) GetResponseStream(result chan<- string, m *Messages) (*openai.ChatCompletionAccumulator, error) {
	stream := lc.Client.Chat.Completions.NewStreaming(context.TODO(), openai.ChatCompletionNewParams{
		Messages: m.ToMessageParamUnion(),
		Model:    lc.Model,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: param.NewOpt(true),
		},
	},
		option.WithJSONSet("enable_thinking", false))
	acc := openai.ChatCompletionAccumulator{}
	defer close(result)
	for stream.Next() {
		chunk := stream.Current()
		//res, _ := json.Marshal(chunk)
		//fmt.Println(string(res))
		acc.AddChunk(chunk)

		//if content, ok := acc.JustFinishedContent(); ok {
		//	println("Content stream finished:", content)
		//}
		//
		//// if using tool calls
		//if tool, ok := acc.JustFinishedToolCall(); ok {
		//	println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
		//}
		//
		//if refusal, ok := acc.JustFinishedRefusal(); ok {
		//	println("Refusal stream finished:", refusal)
		//}

		// it's best to use chunks after handling JustFinished events
		if len(chunk.Choices) > 0 {
			r := chunk.Choices[0].Delta.Content
			//println(r)
			result <- r
		}

	}

	if stream.Err() != nil {
		return nil, stream.Err()
	}
	fmt.Println(acc.ChatCompletion)
	return &acc, nil
}
