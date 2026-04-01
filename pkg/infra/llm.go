package infra

import (
	"cyblog/pkg/llm"

	"github.com/spf13/viper"
)

var globalLLMClient *llm.LLMClient

func NewLLmClient(vc *viper.Viper) *llm.LLMClient {
	globalLLMClient = llm.NewLLMClient(vc)
	return globalLLMClient
}

func GetLLMClient() *llm.LLMClient {
	return globalLLMClient
}
