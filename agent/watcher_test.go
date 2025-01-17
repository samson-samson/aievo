package agent

import (
	"context"
	"fmt"
	"testing"

	"github.com/antgroup/aievo/schema"
)

func TestWatcher(t *testing.T) {

	watcher, err := NewWatcherAgent(WithLLM(client),
		WithPrompt(`You are now in the Werewolf game. Your goal is to sort out which player is out at the end of each round of the game and eliminate this player from the game.`),
		WithInstruction(`
Survival players in the game:
~~~
{{.agent_descriptions}}
~~~

When you need to remove players from game, your Answer must be json format like:
{
    "remove": ["AGENT NAME", "AGENT NAME"],
}

When there is no player to be remove from game, your Answer must be json format like:
{
    "remove": [],
}

`),
		WithSuffix(`
Game conversation history:
{{.history}}

Now it is your turn to give your answer, Begin!

Answer:`))

	if err != nil {
		panic(err)
	}
	run, err := watcher.Run(context.Background(),
		[]schema.Message{
			{
				Receiver:  "ALL",
				Type:      "msg",
				Thought:   "宣布游戏开始并介绍规则，确保所有玩家了解游戏流程和各自的角色职责。",
				Content:   "欢迎来到狼人杀游戏！游戏中有3名狼人、3名神职和3名村民。游戏分为夜晚和白天两个阶段。夜晚阶段，狼人、神职依次行动；白天阶段，所有玩家讨论并投票淘汰一名玩家。游戏将持续进行，直到一方阵营获胜。现在，游戏正式开始！",
				Condition: "daily",
			}})
	if err != nil {
		panic(err)
	}
	fmt.Println(run)
}
