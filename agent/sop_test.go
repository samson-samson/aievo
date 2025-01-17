package agent

import (
	"context"
	"testing"

	"github.com/antgroup/aievo/environment"
	"github.com/antgroup/aievo/feedback"
	"github.com/antgroup/aievo/schema"
)

func TestSOP(t *testing.T) {

	sopFeedback, _ := feedback.NewSopFeedback(client())
	sop, err := NewSopAgent(WithLLM(client()),
		// agent.WithTools([]tool.Tool{sr}),
		WithEnv(environment.NewEnv()),
		WithFeedbacks(sopFeedback))
	if err != nil {
		t.Fatal(err)
	}
	result, err := sop.Run(context.Background(), []schema.Message{
		{
			Type:    schema.MsgTypeMsg,
			Content: "开始一个狼人杀游戏，有3个狼，3个神，3个村民，角色已经分配，请制定游戏的SOP，并开始游戏",
		}})
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(result)
}
