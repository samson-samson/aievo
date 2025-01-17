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

type DeleteFolder struct {
	Workspace string
}

var _ tool.Tool = &DeleteFolder{}

func NewDeleteFolder(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &DeleteFolder{
		Workspace: options.Workspace,
	}, nil
}

func (t *DeleteFolder) Name() string {
	return "DeleteFolder"
}

func (t *DeleteFolder) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Delete a folder in the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"foldername\": \"example_folder\"}`
}

func (t *DeleteFolder) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"foldername": {
				Type:        tool.TypeString,
				Description: "The name of the folder to delete",
			},
		},
		Required: []string{"foldername"},
	}
}

func (t *DeleteFolder) Strict() bool {
	return true
}

func (t *DeleteFolder) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Foldername string `json:"foldername"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	folderPath := filepath.Join(t.Workspace, param.Foldername)
	err = os.RemoveAll(folderPath)
	if err != nil {
		return fmt.Sprintf("failed to delete folder: %v", err), nil
	}

	return fmt.Sprintf("Folder deleted successfully at %s", folderPath), nil
}
