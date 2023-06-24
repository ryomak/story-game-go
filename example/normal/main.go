package main

import (
	"context"
	"fmt"
	game "github.com/ryomak/story-game-go"
	"github.com/sashabaranov/go-openai"
)

func main() {

	game, err := game.NewGameBuilder().
		SetTitle("ゼルダの伝説").
		SetInitialStory(`主人公はハイリア王国の騎士。ハイリア王国は魔王のガノンドロフの復活によって乗っ取られてしまった。主人公は力を蓄えながら、ガノンドロフと倒すための旅に出る`).
		SetMe(&game.Character{
			Name:        "リンク",
			MeRelation:  "ユーザ本人",
			Personality: "人間。とてつもない勇気を持っている。騎士",
			Context:     "ユーザが操作するキャラクター。主人公",
			Variables: map[string]any{
				"体力":  100,
				"レベル": 1,
			},
		}).
		AddConstraint("ゲームマスターとしてロールプレイし、基本的な情報を提供し、ユーザーの入力を待ちます。").
		AddConstraint("ゲームの進行状況やキャラクターのステータスをリストで表示する").
		AddConstraint("ユーザが選択できる選択肢を4つ提示してください。また。ユーザ独自の入力もすべて受け付けてください").
		AddConstraint("ユーザの入力をもとに世界を改変してください").
		AddCharacter(&game.Character{
			Name:        "ゼルダ",
			MeRelation:  "国王の娘",
			Personality: "主人公が守る国王の娘",
			Context:     "ユーザとともに旅をするキャラクター。時を操る力をもつ",
			Variables: map[string]any{
				"体力": 100,
				"運":  200,
			},
		}).
		SetGoal("プレイヤーが死亡するか、ガノンドロフが死亡すると終了する").
		SetDifficult("難しい").
		SetChatGPTModel(openai.GPT40613).
		SetIsFunctionCall(false).
		Build()
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()
	if err := game.Start(ctx); err != nil {
		panic(err)
	}
	displayGameLatestData(game)

	for {
		// 標準入力待ち
		fmt.Print(">> ")
		var s string
		if _, err := fmt.Scan(&s); err != nil {
			continue
		}
		if err := game.UserInput(ctx, s); err != nil {
			fmt.Println("入力エラー", err.Error())
			continue
		}
		displayGameLatestData(game)
		if game.IsDone {
			fmt.Println("ゲーム終了")
			break
		}
	}

}

func displayGameLatestData(game *game.Game) {
	stories := game.GetStories()
	if len(stories) == 0 {
		return
	}
	fmt.Println(stories[len(stories)-1].Result.ToText())
	for k, v := range stories[len(stories)-1].Properties {
		fmt.Println(k.Name, v)
	}
}
