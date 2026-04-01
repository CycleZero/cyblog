package llm

import (
	"github.com/openai/openai-go/v3"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string
}
type Messages struct {
	Data []Message
}

func (m *Messages) ToMessageParamUnion() []openai.ChatCompletionMessageParamUnion {
	var u []openai.ChatCompletionMessageParamUnion
	for _, v := range m.Data {
		if v.Role == "user" {
			u = append(u, openai.UserMessage(v.Content))
		} else if v.Role == "assistant" {
			u = append(u, openai.AssistantMessage(v.Content))
		} else if v.Role == "system" {
			u = append(u, openai.SystemMessage(v.Content))
		}
	}
	return u
}
