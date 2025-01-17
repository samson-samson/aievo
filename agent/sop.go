package agent

import (
	"errors"
	"regexp"
	"strings"

	"github.com/antgroup/aievo/schema"
	"github.com/antgroup/aievo/utils/json"
	"github.com/goccy/go-graphviz"
)

var (
	_dotParse = regexp.MustCompile("(?s)```dot\n(.*?)\n```")
)

func NewSopAgent(opts ...Option) (schema.Agent, error) {
	opts = append(defaultSopOptions(), opts...)

	return NewBaseAgent(opts...)
}

func parseSopOutput(_, content string) ([]schema.StepAction, []schema.Message, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, nil, errors.New("content is empty")
	}
	compile := regexp.MustCompile(_jsonParse)
	submatch := compile.FindAllStringSubmatch(content, -1)
	if len(submatch) != 0 {
		content = strings.TrimSpace(submatch[0][1])
	}
	submatch = _dotParse.FindAllStringSubmatch(content, -1)
	dot := content
	if len(submatch) != 0 {
		dot = strings.TrimSpace(submatch[0][1])
	}

	// parse action
	action := &schema.StepAction{}
	err := json.Unmarshal([]byte(dot), &action)
	if err == nil && action.Action != "" {
		return []schema.StepAction{*action}, nil, nil
	}

	if _, err = graphviz.ParseBytes([]byte(dot)); err != nil {
		return nil, nil, errors.New("parse sop to dot graph failed, err: " + err.Error())
	}

	finish := schema.Message{
		Type:    schema.MsgTypeSOP,
		Content: dot,
		Thought: content,
		Log:     content,
	}
	return nil, []schema.Message{finish}, nil
}
