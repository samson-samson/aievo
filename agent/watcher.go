package agent

import (
	"regexp"
	"strings"

	"github.com/antgroup/aievo/schema"
	"github.com/antgroup/aievo/utils/json"
)

func NewWatcherAgent(opts ...Option) (schema.Agent, error) {
	opts = append(defaultWatcherOptions(), opts...)

	return NewBaseAgent(opts...)
}

func parseMngInfoOutput(_, content string) ([]schema.StepAction, []schema.Message, error) {
	content = strings.TrimSpace(content)
	compile := regexp.MustCompile(_jsonParse)
	submatch := compile.FindAllStringSubmatch(content, -1)
	if len(submatch) != 0 {
		content = submatch[0][1]
	}

	// parse action
	action := schema.StepAction{}
	err := json.Unmarshal([]byte(content), &action)
	if err != nil {
		return nil, nil, err
	}
	if action.Action != "" {
		return []schema.StepAction{action}, nil, nil
	}

	finish := schema.Message{
		Type:    schema.MsgTypeCreative,
		MngInfo: &schema.MngInfo{},
	}
	err = json.Unmarshal([]byte(content), finish.MngInfo)
	return nil, []schema.Message{finish}, err
}
