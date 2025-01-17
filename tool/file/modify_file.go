package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

type ModifyFile struct {
	Workspace string
}

var _ tool.Tool = &ModifyFile{}

func NewModifyFile(opts ...Option) (tool.Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Workspace == "" {
		return nil, errors.New("workspace is not specified")
	}

	return &ModifyFile{
		Workspace: options.Workspace,
	}, nil
}

func (t *ModifyFile) Name() string {
	return "ModifyFile"
}

func (t *ModifyFile) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Modify specific lines in a file within the specified workspace, the input must be json schema: %s", string(bytes)) + `
Example Input: {\"filename\": \"example.txt\", \"modifications\": {\"3\": \"New content for line 3\", \"5\": \"New content for line 5\"}}`
}

func (t *ModifyFile) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"filename": {
				Type:        tool.TypeString,
				Description: "The name of the file to modify",
			},
			"modifications": {
				Type:        tool.TypeJson,
				Description: "A map of line numbers (as strings, 1-based index) to new content for those lines",
			},
		},
		Required: []string{"filename", "modifications"},
	}
}

func (t *ModifyFile) Strict() bool {
	return true
}

func (t *ModifyFile) Call(_ context.Context, input string) (string, error) {
	var param struct {
		Filename      string            `json:"filename"`
		Modifications map[string]string `json:"modifications"`
	}

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	filePath := filepath.Join(t.Workspace, param.Filename)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("failed to read file: %v", err), nil
	}

	lines := strings.Split(string(fileContent), "\n")
	for lineNumStr, newContent := range param.Modifications {
		lineNum, err := strconv.Atoi(lineNumStr)
		if err != nil {
			return fmt.Sprintf("invalid line number: %s", lineNumStr), nil
		}
		if lineNum < 1 || lineNum > len(lines) {
			return fmt.Sprintf("line number %d is out of range", lineNum), nil
		}
		lines[lineNum-1] = newContent
	}

	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Sprintf("failed to write file: %v", err), nil
	}

	return fmt.Sprintf("File %s modified successfully at lines: %v", filePath, getSortedKeys(param.Modifications)), nil
}

// Helper function to get sorted keys from a map[string]string
func getSortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort keys (optional, for consistent output)
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}
