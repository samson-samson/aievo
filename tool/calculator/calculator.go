package calculator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/antgroup/aievo/tool"
	"go.starlark.net/lib/math"
	"go.starlark.net/starlark"
)

// Calculator is a tool that can do math.
type Calculator struct{}

var _ tool.Tool = Calculator{}

// Description returns a string describing the calculator tool.
func (c Calculator) Description() string {
	schema, _ := json.Marshal(c.Schema())
	return `Useful for getting the result of a math expression. this is input schema for this tool: 
` + string(schema)
}

// Name returns the name of the tool.
func (c Calculator) Name() string {
	return "calculator"
}

func (c Calculator) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"param": {
				Type:        tool.TypeString,
				Description: "The param should be a valid mathematical expression that could be executed by a starlark evaluator.",
			},
		},
		Required: []string{"param"},
	}
}

func (c Calculator) Strict() bool {
	return true
}

// Call evaluates the input using a starlak evaluator and returns the result as a
// string. If the evaluator errors the error is given in the result to give the
// agent the ability to retry.
func (c Calculator) Call(ctx context.Context, input string) (string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(input), &m)
	if err != nil {
		return "", err
	}
	v, err := starlark.Eval(&starlark.Thread{Name: "main"}, "input", m["param"], math.Module.Members)
	if err != nil {
		return fmt.Sprintf("error from evaluator: %s", err.Error()), nil //nolint:nilerr
	}
	result := v.String()
	return result, nil
}
