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

type RenameFolder struct {
	Workspace string
}

var _ tool.Tool = &RenameFolder{}

func NewRenameFolder(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &RenameFolder{
		Workspace: options.Workspace,
	}, nil
}

func (t *RenameFolder) Name() string {
	return "RenameFolder"
}

func (t *RenameFolder) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Rename folders in the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"oldFoldername\": \"old_folder\", \"newFoldername\": \"new_folder\"}`
}

func (t *RenameFolder) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"oldFoldername": {
				Type:        tool.TypeString,
				Description: "The current name of the folder to rename",
			},
			"newFoldername": {
				Type:        tool.TypeString,
				Description: "The new name for the folder",
			},
		},
		Required: []string{"oldFoldername", "newFoldername"},
	}
}

func (t *RenameFolder) Strict() bool {
	return true
}

func (t *RenameFolder) Call(_ context.Context, input string) (string, error) {
	var param struct {
		OldFoldername string `json:"oldFoldername"`
		NewFoldername string `json:"newFoldername"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	oldPath := filepath.Join(t.Workspace, param.OldFoldername)
	newPath := filepath.Join(t.Workspace, param.NewFoldername)

	// Check if the old folder exists
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Sprintf("folder %s does not exist", oldPath), nil
	}

	// Rename the folder
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Sprintf("failed to rename folder: %v", err), nil
	}

	return fmt.Sprintf("Folder renamed successfully from %s to %s", oldPath, newPath), nil
}
