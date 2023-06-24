# story-game-go
## example
```go
package main

import (
	"context"
	"fmt"
	game "github.com/ryomak/story-game-go"
	"github.com/sashabaranov/go-openai"
	"os"
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

```

output: 
```go
ゲームを開始するため、まずはゲームの現状から説明します。

ールプレイの舞台は、魔王ガノンドロフによって乗っ取られたハイリア王国。あなた、リンクは勇猛な騎士で、まだレベル1、体力は100です。そしてあなたのパートナー、王女ゼルダは体力100、運200を持っています。サイドキックとしてあなたを旅をサポートします。

目の前にはダークフォレストが広がっています。魔物たちが跋扈し、ところどころに触れると痛みを与えるトゲトゲした草木が生えています。

その先には、伝説の剣マスターソードが眠っていると言われています。

さて、リンク。あなたの行動は何ですか？

1. ゼルダと一緒にダークフォレストに突入する
2. 近くの村で情報を集める
3. まずは訓練場で鍛錬を積む
4. 旅の準備をし直す

<カスタムの選択肢も受け付けます>
>> ココカラ村で家を建てる
その選択は少し予想外でしたが、リンクが決めたのならそれに従いましょう。ココカラ村へ向かっていきます。

あなたがココカラ村で立派な家を建てることで、住民たちがあなたを一目見るや、敬意を表し、おそらくあなたは彼らから様々な情報と援助を得られるでしょう。しかし、それは一定の資源と時間を必要とします。建築したい家のデザインやサイズによりますが、その難易度は上昇します。

さあさ、今すぐにでも家作りを始めますか？それとも他の選択を選びますか？

1. ココカラ村で家づくりを始める
2. 他の村民と交流し、情報を探る
3. 近くの森で材料を探す
4. ゼルダと一緒にダンジョンを探索する

<カスタムの選択肢も受け付けます>
>> 1
素晴らしい、あなたの家作りが決定しました！今からあなたの新たなる冒険、家作りが開始します。

まずは、時間とリソースを管理しながら、ツールや資材を見つけ、村の周辺で建築用木材を集め、建設現場を確保する必要があります。このプロセスは容易ではないかもしれませんが、あなたの勇気と努力で必ず成功します。

ゼルダとあなたが各々のタスクに励む様子を見て、村人たちは少しずつ心を開いていきます。そこから新たな情報やヘルプを得られるかもしれません。

さあ、ハンマーや木材を手にとって、家作りを始めましょう！あなたの冒険はこれからだ。

＜数日後＞

あなたとゼルダの努力によって、素晴らしい家がココカラ村に完成しました！下記のステータスでゲームを進行していきます。

- リンクの現在の状態
  - レベル：1 → 2（家作りの経験によりレベルアップ）
  - 体力：100 → 90（家作りの労働により少し疲れました）
- ゼルダの現在の状態
  - 体力：100 → 90（リンクをサポートするために少し疲れました）
  - 運：200

次に何をしますか？

1. 家で一休みする
2. 村の人々に話しかけて情報を集める
3. まだ日があるので森でリソースを集める
4. ゼルダと計画を立てる

<カスタムの選択肢も受け付けます>
>>
```