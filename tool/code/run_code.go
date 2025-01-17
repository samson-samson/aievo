package code

import (
	"context"
	"fmt"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

const (
	_defaultLangType = "golang"
)

var RunnerFactories = map[string]RunnerFactory{
	"golang": NewGolangRunner,
	"python": NewPythonRunner,
	"java":   NewJavaRunner,
}

type Tool struct {
	ProgramLangType string
	runner          Runner
}

var _ tool.Tool = &Tool{}

// New creates a new Tool instance
func New(opts ...Option) (*Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.ProgramLangType == "" {
		options.ProgramLangType = _defaultLangType
	}

	return &Tool{
		ProgramLangType: options.ProgramLangType,
		runner:          RunnerFactories[options.ProgramLangType](),
	}, nil
}

// Name returns the name of the tool
func (t *Tool) Name() string {
	return fmt.Sprintf("%sRunner", t.ProgramLangType)
}

// Description returns the description of the tool
func (t *Tool) Description() string {
	return fmt.Sprintf("Run %s code, the input is %s code",
		t.ProgramLangType, t.ProgramLangType)
}

func (t *Tool) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"code": {
				Type:        tool.TypeString,
				Description: fmt.Sprintf("The %s code to run.", t.ProgramLangType),
			},
		},
		Required: []string{"code"},
	}
}

func (t *Tool) Strict() bool {
	return true
}

// Call runs the tool with the given input
func (t *Tool) Call(_ context.Context, input string) (string, error) {
	var m map[string]interface{}

	err := json.Unmarshal([]byte(input), &m)
	if err != nil {
		return "json unmarshal error, please try agent", nil
	}

	if m["code"] == nil {
		return "code is required", nil
	}

	// check program runtime
	err = t.runner.CheckRuntime()
	if err != nil {
		return err.Error(), nil
	}

	// run code
	result, err := t.runner.Run(m["code"].(string), nil)
	if err != nil {
		return err.Error(), nil
	}

	return result, nil
}
