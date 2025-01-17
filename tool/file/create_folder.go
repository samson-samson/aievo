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

type CreateFolder struct {
	Workspace string
}

var _ tool.Tool = &CreateFolder{}

func NewCreateFolder(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &CreateFolder{
		Workspace: options.Workspace,
	}, nil
}

func (t *CreateFolder) Name() string {
	return "CreateFolder"
}

func (t *CreateFolder) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Create folders in the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"foldername\": \"example_folder\"}`
}

func (t *CreateFolder) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"foldername": {
				Type:        tool.TypeString,
				Description: "The name of the folder to create",
			},
		},
		Required: []string{"foldername"},
	}
}

func (t *CreateFolder) Strict() bool {
	return true
}

func (t *CreateFolder) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Foldername string `json:"foldername"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	folderPath := filepath.Join(t.Workspace, param.Foldername)
	err = os.MkdirAll(folderPath, 0755)
	if err != nil {
		return fmt.Sprintf("failed to create folder: %v", err), nil
	}

	return fmt.Sprintf("Folder created successfully at %s", folderPath), nil
}
