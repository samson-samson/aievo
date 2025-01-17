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

type RenameFile struct {
	Workspace string
}

var _ tool.Tool = &RenameFile{}

func NewRenameFile(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &RenameFile{
		Workspace: options.Workspace,
	}, nil
}

func (t *RenameFile) Name() string {
	return "RenameFile"
}

func (t *RenameFile) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Rename files in the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"oldname\": \"old.txt\", \"newname\": \"new.txt\"}`
}

func (t *RenameFile) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"oldname": {
				Type:        tool.TypeString,
				Description: "The current name of the file to rename",
			},
			"newname": {
				Type:        tool.TypeString,
				Description: "The new name for the file",
			},
		},
		Required: []string{"oldname", "newname"},
	}
}

func (t *RenameFile) Strict() bool {
	return true
}

func (t *RenameFile) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Oldname string `json:"oldname"`
		Newname string `json:"newname"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	oldPath := filepath.Join(t.Workspace, param.Oldname)
	newPath := filepath.Join(t.Workspace, param.Newname)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Sprintf("failed to rename file: %v", err), nil
	}

	return fmt.Sprintf("File renamed successfully from %s to %s", oldPath, newPath), nil
}
