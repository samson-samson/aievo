package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

type CreateFile struct {
	Workspace string
}

var _ tool.Tool = &CreateFile{}

func NewCreateFile(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &CreateFile{
		Workspace: options.Workspace,
	}, nil
}

func (t *CreateFile) Name() string {
	return "CreateFile"
}

func (t *CreateFile) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Create files in the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {"filename": "example.txt", "content": "This is an example file."}`
}

func (t *CreateFile) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"filename": {
				Type:        tool.TypeString,
				Description: "The name of the file to create",
			},
			"content": {
				Type:        tool.TypeString,
				Description: "The content to write to the file",
			},
		},
		Required: []string{"filename", "content"},
	}
}

func (t *CreateFile) Strict() bool {
	return true
}

func (t *CreateFile) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	filePath := filepath.Join(t.Workspace, param.Filename)
	err = os.WriteFile(filePath, []byte(param.Content), 0644)
	if err != nil {
		return fmt.Sprintf("failed to create file: %v", err), nil
	}

	return fmt.Sprintf("File created successfully at %s", filePath), nil
}
