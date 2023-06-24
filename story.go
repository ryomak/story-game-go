package game

import (
	"encoding/json"
	"github.com/sashabaranov/go-openai"
)

type Result interface {
	ToMessage() openai.ChatCompletionMessage
	ToText() string
}

type UserPhase struct {
	Text string
}

func (u *UserPhase) ToMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: u.Text,
	}
}

func (u *UserPhase) ToText() string {
	return u.Text
}

type GMPhase struct {
	Text         string
	FunctionCall *openai.FunctionCall
}

func (g *GMPhase) ToMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:         openai.ChatMessageRoleAssistant,
		Content:      g.Text,
		FunctionCall: g.FunctionCall,
	}
}

func (g *GMPhase) ToText() string {
	if g.FunctionCall != nil {
		properties := map[string]any{}
		_ = json.Unmarshal([]byte(g.FunctionCall.Arguments), &properties)
		res, ok := properties[PropertyKeyGameState]
		if ok {
			return res.(string)
		}
	}
	return g.Text
}
