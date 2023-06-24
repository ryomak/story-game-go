package game

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sashabaranov/go-openai"
)

const (
	PropertyKeyIsDone    = "isDone"
	PropertyKeyGameState = "gameMasterStatement"

	FunctionCallNameGetStoryState = "get_story_state"
)

type Game struct {
	story           []*Story
	parameterMap    map[ParameterKey]any
	parameterKeyMap map[string]ParameterKey
	IsDone          bool

	title        string
	goal         string
	difficult    string
	initialStory string
	constraints  []string
	me           *Character
	characterMap map[string]*Character

	openAIClient   *openai.Client
	chatGPTModel   string
	isFunctionCall bool
}

type Story struct {
	Result Result

	Properties map[ParameterKey]any
}

func (g *Game) Start(ctx context.Context) error {
	g.story = append(g.story, &Story{
		Result: &UserPhase{
			Text: "ゲームを開始します",
		},
	})
	res, err := g.openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:        g.chatGPTModel,
		Messages:     g.chatGPTMessages(),
		MaxTokens:    1000,
		TopP:         1,
		Temperature:  0,
		Functions:    g.functionCallDefinition(),
		FunctionCall: g.functionCall(),
	})
	if err != nil {
		return err
	}
	story, err := g.parseStory(res)
	if err != nil {
		return err
	}
	g.story = append(g.story, story)
	return nil
}

func (g *Game) UserInput(ctx context.Context, s string) error {
	g.story = append(g.story, &Story{
		Result: &UserPhase{
			Text: s,
		},
	})
	res, err := g.openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:        g.chatGPTModel,
		Messages:     g.chatGPTMessages(),
		MaxTokens:    1000,
		TopP:         1,
		Temperature:  0,
		Functions:    g.functionCallDefinition(),
		FunctionCall: g.functionCall(),
	})
	if err != nil {
		return err
	}
	story, err := g.parseStory(res)
	if err != nil {
		return err
	}
	if isDone, ok := story.Properties[g.parameterKeyMap[PropertyKeyIsDone]].(bool); isDone && ok {
		g.IsDone = true
	}
	g.story = append(g.story, story)
	return nil
}

func (g *Game) GetStories() []*Story {
	return g.story
}

func (g *Game) parseStory(res openai.ChatCompletionResponse) (*Story, error) {
	if len(res.Choices) == 0 {
		return nil, errors.New("no choices")
	}
	properties := make(map[ParameterKey]any)

	choice := res.Choices[0]
	if call := choice.Message.FunctionCall; call != nil {
		r := make(map[string]any)
		if err := json.Unmarshal([]byte(call.Arguments), &r); err != nil {
			return nil, err
		}
		for k, v := range r {
			key, ok := g.parameterKeyMap[k]
			if !ok {
				continue
			}
			properties[key] = v
		}
	}
	return &Story{
		Result: &GMPhase{
			Text:         choice.Message.Content,
			FunctionCall: choice.Message.FunctionCall,
		},
		Properties: properties,
	}, nil
}

func (g *Game) chatGPTMessages() []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0, len(g.story)+1)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: g.systemPrompt(),
	})
	for _, v := range g.story {
		messages = append(messages, v.Result.ToMessage())
	}
	return messages
}

type ParameterType string

const (
	PropertyTypeString     ParameterType = "string"
	PropertyTypeNumber     ParameterType = "number"
	PropertyTypeBoolean    ParameterType = "boolean"
	PropertyTypeListString ParameterType = "listString"
	PropertyTypeListNumber ParameterType = "listNumber"
)

type ParameterKey struct {
	Name        string
	Description string
	Type        ParameterType
}

func (g *Game) functionCall() any {
	if g.isFunctionCall {
		return map[string]any{"name": FunctionCallNameGetStoryState}
	}
	return "none"
}

func (g *Game) functionCallDefinition() []openai.FunctionDefinition {
	requires := make([]string, 0, len(g.parameterMap))
	properties := make(map[string]openai.JSONSchemaDefinition)
	for k := range g.parameterMap {
		switch k.Type {
		case PropertyTypeString:
			properties[k.Name] = openai.JSONSchemaDefinition{
				Type:        openai.JSONSchemaTypeString,
				Description: k.Description,
			}
		case PropertyTypeNumber:
			properties[k.Name] = openai.JSONSchemaDefinition{
				Type:        openai.JSONSchemaTypeNumber,
				Description: k.Description,
			}
		case PropertyTypeBoolean:
			properties[k.Name] = openai.JSONSchemaDefinition{
				Type:        openai.JSONSchemaTypeBoolean,
				Description: k.Description,
			}
		case PropertyTypeListString:
			properties[k.Name] = openai.JSONSchemaDefinition{
				Type: openai.JSONSchemaTypeArray,
				Items: &openai.JSONSchemaDefinition{
					Type: openai.JSONSchemaTypeString,
				},
				Description: k.Description,
			}
		case PropertyTypeListNumber:
			properties[k.Name] = openai.JSONSchemaDefinition{
				Type: openai.JSONSchemaTypeArray,
				Items: &openai.JSONSchemaDefinition{
					Type: openai.JSONSchemaTypeNumber,
				},
				Description: k.Description,
			}
		default:
			continue
		}
		requires = append(requires, k.Name)
	}
	return []openai.FunctionDefinition{
		{
			Name:        "get_story_state",
			Description: "ゲームの状態を取得する",
			Parameters: &openai.JSONSchemaDefinition{
				Type:       openai.JSONSchemaTypeObject,
				Properties: properties,
				Required:   requires,
			},
		},
	}
}
