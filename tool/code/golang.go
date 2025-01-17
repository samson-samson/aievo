package code

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GolangRunner struct {
}

// NewGolangRunner creates a new instance of GolangRunner
func NewGolangRunner() Runner {
	return &GolangRunner{}
}

// CheckRuntime checks if the Go runtime environment is installed and accessible
func (r *GolangRunner) CheckRuntime() error {
	// Try to execute "go version" to check if Go is installed
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()

	if err != nil {
		// If there's an error executing the command, Go is probably not installed
		return errors.New("go runtime environment not found: " + err.Error())
	}

	// Check if the output contains the word "go version"
	if !strings.Contains(string(output), "go version") {
		return errors.New("go runtime environment is not correctly installed")
	}

	// If everything is fine, return nil (no error)
	return nil
}

// Run compiles and executes the Go code with the provided arguments using "go run"
func (r *GolangRunner) Run(code string, args []string) (string, error) {
	// Create a temporary directory for the Go files
	tmpDir, err := os.MkdirTemp("", "golang_runner")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	// Create the Go file inside the temporary directory
	goFilePath := filepath.Join(tmpDir, "main.go")
	goFile, err := os.Create(goFilePath)
	if err != nil {
		return "", err
	}

	// Write the Go code to the temporary file
	if _, err := goFile.Write([]byte(code)); err != nil {
		return "", err
	}
	if err := goFile.Close(); err != nil {
		return "", err
	}

	// Run the Go code using "go run" and pass arguments
	cmd := exec.Command("go", append([]string{"run", goFilePath}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), errors.New("error during execution: " + err.Error())
	}

	return string(output), nil
}
