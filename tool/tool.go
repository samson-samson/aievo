package tool

import "context"

var (
	TypeJson   = "object"
	TypeString = "string"
	TypeInt    = "integer"
	TypeArr    = "array"
)

type PropertiesSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]PropertySchema `json:"properties,omitempty"`
	Required   []string                  `json:"required"`
	// for azure
	AdditionalProperties bool `json:"additionalProperties"`
}
type PropertySchema struct {
	Type        string                    `json:"type"`
	Description string                    `json:"description"`
	Enum        []string                  `json:"enum,omitempty"`
	Properties  map[string]PropertySchema `json:"properties,omitempty"`
	Required    []string                  `json:"required,omitempty"`

	AdditionalProperties bool               `json:"additionalProperties"`
	Default              string             `json:"default,omitempty"`
	Items                *PropertySchema    `json:"items,omitempty"`
	OneOf                []PropertiesSchema `json:"oneOf,omitempty"`
}

type Tool interface {
	Name() string
	Description() string
	Schema() *PropertiesSchema
	Strict() bool
	Call(context.Context, string) (string, error)
}
