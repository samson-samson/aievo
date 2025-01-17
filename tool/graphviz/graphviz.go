package graphviz

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/antgroup/aievo/llm"
	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

var (
	ErrMissLLM      = errors.New("missing LLM")
	ErrParamBeBlank = errors.New("agent_descriptions/sop cannot be blank")
	ErrExecFailed   = errors.New("exec tool failed, please retry")
	_dotParse       = regexp.MustCompile("(?s)```dot\n(.*?)\n```")
	_graphvizPrompt = `
You are a helpful assistant to help user generate flowchart by using dot language and the following:

Instructions:
1. Need to list node information separately
2. Using the default shapes, the layout uses top-down
3. The judgment logic does not use separate nodes
4. Each node requires Agent participation
5. An Agent can only correspond to one Node, even if the Agent can do different tasks
6. The number of nodes is less than or equal to the number of Agents
7. The format of the node information label is {"agent":"","action":""}
8. The start of the process requires a separate node, this node's label is {"agent":"start","action":"start"}
9. The end of the process requires a separate node, this node's label is {"agent":"end","action":"end"}
10. Each node's status is uncompleted by default
11. If there is more than one out-degree, conditions need to be indicated

Existing agents:
%s

You must response with json with format:` +
		"\n```dot\n" +
		"This is dot content\n" +
		"```" + `

Here is Sop to convert to dot graphviz:
%s

history:
%s

Output:
`
)

type Graphviz struct {
	prompt   string
	llm      llm.LLM
	maxRetry int
}

type GraphvizRequest struct {
	AgentDescriptions string `json:"agent_descriptions"`
	Sop               string `json:"sop"`
}

func NewGraphvizTool(llm llm.LLM, maxRetry int) (*Graphviz, error) {
	if llm == nil {
		return nil, ErrMissLLM
	}
	if maxRetry < 0 {
		maxRetry = 1
	}
	return &Graphviz{
		llm:      llm,
		maxRetry: maxRetry,
		prompt:   _graphvizPrompt,
	}, nil
}

func (*Graphviz) Name() string {
	return "SOP To Dot Graph"
}

func (*Graphviz) Description() string {
	return "This tool integrates with Graphviz to " +
		"automatically convert Standard Operating Procedure (SOP) " +
		"documents into the DOT language, which is used for " +
		"creating structured diagrams. " +
		"By converting SOPs into visual workflows, " +
		"it allows the agent to interpret and follow " +
		"the steps in a clear, logical format. " +
		"The diagrams help the agent navigate tasks more efficiently," +
		" ensuring better compliance with procedural guidelines and " +
		"reducing errors in execution. The tool accepts input in JSON format with the following structure:\n" +
		`{
"agent_descriptions": "agent name and description, like Math: the agent is to Calculator math problem\nTranslate: the agent is to translate",
"sop": "sop(Standard Operating Procedure)"
}`
}

func (g *Graphviz) Schema() *tool.PropertiesSchema {
	return nil
}

func (g *Graphviz) Strict() bool {
	return true
}

func (g *Graphviz) Call(ctx context.Context, s string) (string, error) {
	req := &GraphvizRequest{}
	err := json.Unmarshal([]byte(s), req)
	if err != nil {
		return "", err
	}
	if req.AgentDescriptions == "" || req.Sop == "" {
		return "", ErrParamBeBlank
	}
	history := ""
	for i := 0; i < g.maxRetry; i++ {
		prompt := fmt.Sprintf(g.prompt, req.AgentDescriptions,
			req.Sop, history)
		result, err := g.llm.Generate(ctx, prompt)
		if err != nil {
			return "", ErrExecFailed
		}
		_, err = g.parseDot(result.Content)
		if err != nil {
			history += "AI: " + result.Content + "\n" +
				"Feedback: " + err.Error() + "\n"
			continue
		}
		return result.Content, nil
	}
	return "", ErrExecFailed
}

func (g *Graphviz) parseDot(output string) (*cgraph.Graph, error) {
	match := _dotParse.FindStringSubmatch(output)
	if len(match) >= 2 {
		output = match[1]
	}
	return graphviz.ParseBytes([]byte(output))
}
