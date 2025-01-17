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

type ListFolder struct {
	Workspace string
}

var _ tool.Tool = &ListFolder{}

func NewListFolder(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &ListFolder{
		Workspace: options.Workspace,
	}, nil
}

func (t *ListFolder) Name() string {
	return "ListFolder"
}

func (t *ListFolder) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("List files in the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"foldername\": \"example_folder\"}`
}

func (t *ListFolder) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"foldername": {
				Type:        tool.TypeString,
				Description: "The name of the folder to list files from",
			},
		},
		Required: []string{"foldername"},
	}
}

func (t *ListFolder) Strict() bool {
	return true
}

func (t *ListFolder) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Foldername string `json:"foldername"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	folderPath := filepath.Join(t.Workspace, param.Foldername)
	tree, err := t.generateTree(folderPath, "")
	if err != nil {
		return fmt.Sprintf("failed to list files: %v", err), nil
	}

	return tree, nil
}

func (t *ListFolder) generateTree(rootPath string, prefix string) (string, error) {
	var tree strings.Builder

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return "", err
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1
		tree.WriteString(prefix)
		if isLast {
			tree.WriteString("└── ")
		} else {
			tree.WriteString("├── ")
		}
		tree.WriteString(entry.Name())
		tree.WriteString("\n")

		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			subTree, err := t.generateTree(filepath.Join(rootPath, entry.Name()), newPrefix)
			if err != nil {
				return "", err
			}
			tree.WriteString(subTree)
		}
	}

	return tree.String(), nil
}
