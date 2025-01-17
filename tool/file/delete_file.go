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

type DeleteFile struct {
	Workspace string
}

var _ tool.Tool = &DeleteFile{}

func NewDeleteFile(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &DeleteFile{
		Workspace: options.Workspace,
	}, nil
}

func (t *DeleteFile) Name() string {
	return "DeleteFile"
}

func (t *DeleteFile) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Delete files in the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"filename\": \"example.txt\"}`
}

func (t *DeleteFile) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"filename": {
				Type:        tool.TypeString,
				Description: "The name of the file to delete",
			},
		},
		Required: []string{"filename"},
	}
}

func (t *DeleteFile) Strict() bool {
	return true
}

func (t *DeleteFile) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Filename string `json:"filename"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	filePath := filepath.Join(t.Workspace, param.Filename)
	err = os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Sprintf("file does not exist: %s", filePath), nil
		}
		return fmt.Sprintf("failed to delete file: %v", err), nil
	}

	return fmt.Sprintf("File deleted successfully: %s", filePath), nil
}
