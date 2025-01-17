package code

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type JavaRunner struct {
}

// NewJavaRunner creates a new instance of JavaRunner
func NewJavaRunner() Runner {
	return &JavaRunner{}
}

// CheckRuntime checks if the Java runtime environment is installed and accessible
func (r *JavaRunner) CheckRuntime() error {
	// Try to execute "java -version" to check if Java is installed
	cmd := exec.Command("java", "-version")
	output, err := cmd.CombinedOutput() // java -version outputs to stderr

	if err != nil {
		// If there's an error executing the command, Java is probably not installed
		return errors.New("java runtime environment not found: " + err.Error())
	}

	// Check if the output contains the word "java"
	if !strings.Contains(string(output), "java version") && !strings.Contains(string(output), "openjdk version") {
		return errors.New("java runtime environment is not correctly installed")
	}

	// If everything is fine, return nil (no error)
	return nil
}

// Run compiles and executes the provided Java code with the provided arguments using "javac" and "java"
func (r *JavaRunner) Run(code string, args []string) (string, error) {
	// Extract the class name from the Java code using a regular expression
	className, err := extractClassName(code)
	if err != nil {
		return "", err
	}

	// Create a temporary directory for the Java files
	tmpDir, err := os.MkdirTemp("", "java_runner")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	// Create the Java file in the temporary directory with the correct class name
	javaFilePath := filepath.Join(tmpDir, className+".java")
	javaFile, err := os.Create(javaFilePath)
	if err != nil {
		return "", err
	}

	// Write the Java code to the temporary file
	if _, err := javaFile.Write([]byte(code)); err != nil {
		return "", err
	}
	if err := javaFile.Close(); err != nil {
		return "", err
	}

	// Compile the Java code using "javac" inside the temporary directory
	compileCmd := exec.Command("javac", javaFilePath)
	compileOutput, err := compileCmd.CombinedOutput()
	if err != nil {
		return string(compileOutput), errors.New("error during compilation: " + err.Error())
	}

	// Run the compiled Java class using "java" from the temporary directory
	runCmd := exec.Command("java", "-cp", tmpDir, className)
	runCmd.Args = append(runCmd.Args, args...)
	runOutput, err := runCmd.CombinedOutput()
	if err != nil {
		return string(runOutput), errors.New("error during execution: " + err.Error())
	}

	return string(runOutput), nil
}

// extractClassName extracts the class name from the Java code
func extractClassName(code string) (string, error) {
	// Use a regular expression to find the class name
	re := regexp.MustCompile(`(?m)public\s+class\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := re.FindStringSubmatch(code)
	if len(matches) < 2 {
		return "", errors.New("unable to find a public class in the code")
	}
	return matches[1], nil
}
