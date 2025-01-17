package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

type ReadFile struct {
	Workspace string
}

var _ tool.Tool = &ReadFile{}

func NewReadFile(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &ReadFile{
		Workspace: options.Workspace,
	}, nil
}

func (t *ReadFile) Name() string {
	return "ReadFile"
}

func (t *ReadFile) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Read files from the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"filename\": \"example.txt\"}`
}

func (t *ReadFile) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"filename": {
				Type:        tool.TypeString,
				Description: "The name of the file to read",
			},
		},
		Required: []string{"filename"},
	}
}

func (t *ReadFile) Strict() bool {
	return true
}

func (t *ReadFile) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Filename string `json:"filename"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	filePath := filepath.Join(t.Workspace, param.Filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("failed to read file: %v", err), nil
	}

	// Add ">>>" and line numbers to each line
	lines := string(content)
	numberedLines := ""
	lineNumber := 1
	for _, line := range strings.Split(lines, "\n") {
		numberedLines += fmt.Sprintf(">>>%d %s\n", lineNumber, line)
		lineNumber++
	}

	return numberedLines, nil
}
