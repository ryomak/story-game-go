package game

import "fmt"

type Character struct {
	Name        string
	MeRelation  string
	Personality string
	Context     string
	Variables   map[string]any
}

func (c *Character) PromptString() string {
	params := "[追加パラメータ]"
	for k, v := range c.Variables {
		params += fmt.Sprintf("%s: %v\n", k, v)
	}
	return fmt.Sprintf(`[名前]
%s
[性格]
%s
[ユーザとの関係]
%s
[文脈]
%s
[追加パラメータ]
%s
`, c.Name, c.Personality, c.MeRelation, c.Context, params)
}
