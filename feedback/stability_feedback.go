package feedback

import (
	"context"

	"github.com/antgroup/aievo/llm"
	"github.com/antgroup/aievo/schema"
	"github.com/antgroup/aievo/utils/json"
)

var (
	_stabilityPrompt = `
# Goal
You are a task execution review expert, and you are currently evaluating the steps and execution results of various agents within a team. Your goal is to check whether the execution plan meets the requirements and provide improvement suggestions to ensure that the entire team is steadily advancing and executing the plan according to the SOP (Standard Operating Procedure).

# Divided into action and message, this is the action

# Steps
{{if .actions}}
1. You should fully understand the SOP, comprehend the team's communication history, and try to understand which stage of the SOP the current situation is in.
2. You should check whether {{.agent_name}}'s execution plan meets the following requirements:
2.1. Whether {{.agent_name}}'s thinking and the tools {{.agent_name}} executes are consistent.
2.2. When {{.agent_name}} executes the tools, whether the input information is correct and whether there are any hallucinations.
2.3. When {{.agent_name}} executes the tools, whether the scope of the input query is as expected. This point is mainly aimed at monitoring data queries to avoid querying unnecessary metrics.
{{end}}
{{if .messages}}
1. You should fully understand the SOP, comprehend the team's communication history, and try to understand which stage of the SOP the current situation is in.
2. You should check whether {{.agent_name}}'s execution plan meets the following requirements:
2.1. Whether {{.agent_name}}'s thinking and the message {{.agent_name}} outputs are consistent.
2.2. Whether the message {{.agent_name}} outputs is accurate and whether there are any hallucinations.
2.3. Whether the message {{.agent_name}} outputs is complete, including the necessary parameters or data for subsequent execution plans.
2.4. Whether the message {{.agent_name}} outputs contains suggestions. For non-leader agents, no suggestions should be output; they should only process the tasks of other agents and return results.
{{if .leader}}
2.5. Determine whether this message adheres to the SOP.
{{end}}
{{end}}

# Background
## SOP
{{.sop}}

## Tools
here is tools that agent({{.agent_name}}) can select to use:
{{.tool_description}}

## Team Member
here is team members that agent({{.agent_name}}) can communicate with:
{{.agent_description}}

here is tools that agent({{.agent_name}}) can select to use:
{{.tool_description}}

## Team Historical Dialogue
{{.history}}

## Agent History Action
{{.steps}}


# Format example
Your final output should ALWAYS in the following format:
{
	"type": "Approved/NotApproved",
	"msg": "When You are not approved, you should give your suggestion"
}

# Begin
Now start reviewing the following content for agent({{.agent_name}}):
{{.actions}}
{{.messages}}

Output:
`
)

func NewStabilityFeedback(LLM llm.LLM) (Feedback, error) {
	return NewLLMFeedback(LLM,
		WithPromptTemplate(_stabilityPrompt),
		WithMiddlewares(FillInStabilityPrompt))
}

func FillInStabilityPrompt(ctx context.Context, lf *LLMFeedback,
	a schema.Agent, messages []schema.Message,
	actions []schema.StepAction, steps []schema.StepAction,
	_ string) bool {
	action, message, step := []byte(""), []byte(""), []byte("")
	if len(actions) > 0 {
		action, _ = json.Marshal(actions)
	}
	if len(steps) > 0 {
		step, _ = json.Marshal(steps)
	}
	if len(messages) > 0 {
		message, _ = json.Marshal(messages)
	}
	leader := ""
	if leader == a.Env().GetTeamLeader().Name() {
		leader = a.Name()
	}

	// 提取 stability prompt，并提取关键信息
	lf.Vars = map[string]any{
		"actions":           string(action),
		"messages":          string(message),
		"steps":             string(step),
		"leader":            leader,
		"agent_name":        a.Name(),
		"sop":               a.Env().SOP(),
		"tool_descriptions": schema.ConvertToolDescriptions(a.Tools()),
		"agent_descriptions": schema.ConvertAgentDescriptions(
			a.Env().GetSubscribeAgents(ctx, a)),
		"history": schema.ConvertConstructScratchPad(a.Name(), a.Name(),
			messages, steps),
	}
	return true
}
