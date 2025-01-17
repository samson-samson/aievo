package agent

import (
	"context"
	"fmt"
	"testing"

	"github.com/antgroup/aievo/feedback"
	"github.com/antgroup/aievo/llm"
	"github.com/antgroup/aievo/llm/openai"
	"github.com/antgroup/aievo/schema"
	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/tool/calculator"
	"github.com/goccy/go-graphviz"
)

var (
	client, _ = openai.New(openai.WithToken("token"),
		openai.WithModel("auto"),
		openai.WithBaseURL("https://mmcv.alipay.com/api/v1/agent/ant"))
)

func TestBaseAgent(t *testing.T) {
	base, err := NewBaseAgent(
		WithLLM(client),
		WithName("test"),
		WithDesc("test"),
		WithTools([]tool.Tool{
			calculator.Calculator{},
		}),
		WithFeedbacks(&feedback.ContentFeedback{}))
	if err != nil {
		panic(err)
	}
	run, err := base.Run(context.Background(), []schema.Message{
		{
			Sender:   "User",
			Receiver: base.name,
			Content:  "20乘以30等于几",
			Type:     "Msg",
		},
	},
		llm.WithTemperature(0.1),
		llm.WithTopP(0.8),
		llm.WithRepetitionPenalty(1.05),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", run)
}

func TestConversationAgent(t *testing.T) {
	llm, err := openai.New(
		openai.WithModel("auto"),
		openai.WithToken("12345"),
		openai.WithBaseURL("https://mmcv.alipay.com/api/v1/agent"),
	)
	if err != nil {
		panic(err)
	}

	base, err := NewBaseAgent(
		WithLLM(llm),
		WithName("test"),
		WithDesc("test"))
	if err != nil {
		panic(err)
	}
	run, err := base.Run(context.Background(),
		[]schema.Message{
			{
				Sender:   "User",
				Receiver: base.name,
				Content:  "hello, my name is bobby",
				Type:     "Msg",
			},
		})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", run.Messages[0])
	run, err = base.Run(context.Background(),
		[]schema.Message{
			{
				Sender:   "User",
				Receiver: base.name,
				Content:  "what is my name?",
				Type:     "Msg",
			},
		})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", run.Messages[0])
}

func TestGoGraphviz(t *testing.T) {
	path := "test.dot"
	g, err := graphviz.ParseFile(path)
	if err != nil {
		panic(err)
	}
	curNode := g.FirstNode()
	for {
		if curNode == nil {
			break
		}
		fmt.Println(curNode.Name())
		fmt.Println(curNode.Get("label"))
		curNode = g.NextNode(curNode)
	}
}
