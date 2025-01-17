package prompt

import (
	"bytes"
	"fmt"
	"text/template"
)

// PromptTemplate contains common fields for all prompt templates.
type PromptTemplate struct {
	// Template is the prompt template.
	Template *template.Template

	// A list of variable names the prompt template expects.
	InputVariables []string
}

// NewPromptTemplate returns a new prompt template.
func NewPromptTemplate(prompt string) (*PromptTemplate, error) {
	tpl, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return nil, fmt.Errorf("parse prompt template error: %s", err.Error())
	}
	return &PromptTemplate{
		Template: tpl,
	}, nil
}

// Format formats the prompt template and returns a string value.
func (p *PromptTemplate) Format(values map[string]any) (string, error) {
	var buffer bytes.Buffer
	err := p.Template.Execute(&buffer, values)
	if err != nil {
		return "", fmt.Errorf("format prompt template error: %s", err.Error())
	}
	return buffer.String(), nil
}
