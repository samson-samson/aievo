package file

import (
	"fmt"
	"os"

	"github.com/antgroup/aievo/tool"
)

type Factory func(...Option) (tool.Tool, error)

func GetFileRelatedTools(workspace string) ([]tool.Tool, error) {
	if _, err := os.Stat(workspace); os.IsNotExist(err) {
		err := os.MkdirAll(workspace, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create workspace: %v", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check workspace: %v", err)
	}

	return buildFileRelatedTools(
		workspace,
		NewCreateFile,
		NewReadFile,
		NewModifyFile,
		NewDeleteFile,
		NewRenameFile,
		NewCreateFolder,
		NewListFolder,
		NewRenameFolder,
		NewDeleteFolder,
	)
}

func buildFileRelatedTools(workspace string, factories ...Factory) ([]tool.Tool, error) {
	tools := make([]tool.Tool, 0)
	for _, factory := range factories {
		t, err := factory(WithWorkspace(workspace))
		if err != nil {
			return nil, fmt.Errorf("failed to build tool: %v", err)
		}
		tools = append(tools, t)
	}
	return tools, nil
}
