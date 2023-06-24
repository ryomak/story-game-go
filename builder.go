package game

import (
	"errors"
	"github.com/sashabaranov/go-openai"
	"os"
)

type gameBuilder struct {
	title          string
	goal           string
	initialStory   string
	difficult      string
	me             *Character
	characters     []*Character
	openAIClient   *openai.Client
	chatGPTModel   string
	parameters     map[ParameterKey]any
	constraints    []string
	isFunctionCall *bool
}

func NewGameBuilder() *gameBuilder {
	return &gameBuilder{
		characters:  []*Character{},
		constraints: []string{},
	}
}

func (b *gameBuilder) SetTitle(title string) *gameBuilder {
	b.title = title
	return b
}

func (b *gameBuilder) SetGoal(goal string) *gameBuilder {
	b.goal = goal
	return b
}

func (b *gameBuilder) SetInitialStory(story string) *gameBuilder {
	b.initialStory = story
	return b
}

func (b *gameBuilder) SetMe(me *Character) *gameBuilder {
	b.me = me
	return b
}

func (b *gameBuilder) SetOpenAIClient(client *openai.Client) *gameBuilder {
	b.openAIClient = client
	return b
}

func (b *gameBuilder) AddConstraint(constraint string) *gameBuilder {
	b.constraints = append(b.constraints, constraint)
	return b
}

func (b *gameBuilder) AddCharacter(character *Character) *gameBuilder {
	b.characters = append(b.characters, character)
	return b
}

func (b *gameBuilder) SetChatGPTModel(model string) *gameBuilder {
	b.chatGPTModel = model
	return b
}

func (b *gameBuilder) SetDifficult(difficult string) *gameBuilder {
	b.difficult = difficult
	return b
}

func (b *gameBuilder) SetParameters(parameters map[ParameterKey]any) *gameBuilder {
	b.parameters = parameters
	return b
}

func (b *gameBuilder) SetIsFunctionCall(isFunctionCall bool) *gameBuilder {
	b.isFunctionCall = &isFunctionCall
	return b
}

func (b *gameBuilder) Build() (*Game, error) {
	characterMap := map[string]*Character{}
	for _, c := range b.characters {
		characterMap[c.Name] = c
	}
	if b.me == nil {
		return nil, errors.New("me is not set")
	}
	if b.goal == "" {
		return nil, errors.New("goal is not set")
	}
	if b.initialStory == "" {
		return nil, errors.New("initialStory is not set")
	}
	if b.openAIClient == nil {
		b.openAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	}
	if b.chatGPTModel == "" {
		b.chatGPTModel = openai.GPT40613
	}
	if b.title == "" {
		b.title = "RPG"
	}
	if b.difficult == "" {
		b.difficult = "normal"
	}
	if b.parameters == nil {
		b.parameters = map[ParameterKey]any{}
	}
	b.parameters[ParameterKey{
		Name:        PropertyKeyIsDone,
		Description: "ゲームが終了したかどうか",
		Type:        PropertyTypeBoolean,
	}] = false
	b.parameters[ParameterKey{
		Name:        PropertyKeyGameState,
		Description: "ゲームマスターの発言内容。ゲームマスターの発言内容は、ゲームマスターの発言内容のみを含む",
		Type:        PropertyTypeString,
	}] = ""
	parameterKeyMap := map[string]ParameterKey{}
	for k := range b.parameters {
		parameterKeyMap[k.Name] = k
	}
	if b.isFunctionCall == nil {
		return nil, errors.New("isFunctionCall is not set")
	}

	return &Game{
		title:           b.title,
		goal:            b.goal,
		initialStory:    b.initialStory,
		me:              b.me,
		characterMap:    characterMap,
		openAIClient:    b.openAIClient,
		chatGPTModel:    b.chatGPTModel,
		parameterMap:    b.parameters,
		parameterKeyMap: parameterKeyMap,
		constraints:     b.constraints,
		isFunctionCall:  *b.isFunctionCall,
	}, nil
}
