package code

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PythonRunner struct {
}

// NewPythonRunner creates a new instance of PythonRunner
func NewPythonRunner() Runner {
	return &PythonRunner{}
}

// CheckRuntime checks if the Python runtime environment is installed and accessible
func (r *PythonRunner) CheckRuntime() error {
	// Try to execute "python3 --version" to check if Python is installed
	cmd := exec.Command("python3", "--version")
	output, err := cmd.Output()

	if err != nil {
		// If there's an error executing the command, Python is probably not installed
		return errors.New("python runtime environment not found: " + err.Error())
	}

	// Check if the output contains the word "Python"
	if !strings.Contains(string(output), "Python") {
		return errors.New("python runtime environment is not correctly installed")
	}

	// If everything is fine, return nil (no error)
	return nil
}

// Run executes the provided Python code with the provided arguments using "python3"
func (r *PythonRunner) Run(code string, args []string) (string, error) {
	// Create a temporary directory for the Python file
	tmpDir, err := os.MkdirTemp("", "python_runner")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	// Create the Python file inside the temporary directory
	pyFilePath := filepath.Join(tmpDir, "main.py")
	pyFile, err := os.Create(pyFilePath)
	if err != nil {
		return "", err
	}

	// Write the Python code to the temporary file
	if _, err := pyFile.Write([]byte(code)); err != nil {
		return "", err
	}
	if err := pyFile.Close(); err != nil {
		return "", err
	}

	// Run the Python code using "python3" and pass arguments
	cmd := exec.Command("python3", append([]string{pyFilePath}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), errors.New("error during execution: " + err.Error())
	}

	return string(output), nil
}
