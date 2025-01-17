package reader

import (
	"context"
	"errors"
	"fmt"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

var readerFactories = map[string]Factory{
	"website": NewWebsiteReader,
	"pdf":     NewPdfReader,
}

type Tool struct {
	ReaderType string
	Reader     Reader
}

var _ tool.Tool = &Tool{}

func NewReader(opts ...Option) (*Tool, error) {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.ReaderType == "" {
		return nil, errors.New("reader type is not specified")
	}

	readerFactory, ok := readerFactories[options.ReaderType]
	if !ok {
		return nil, errors.New("reader type is not supported")
	}

	return &Tool{
		ReaderType: options.ReaderType,
		Reader:     readerFactory(),
	}, nil

}

func (t *Tool) Name() string {
	return fmt.Sprintf("%sReader", t.ReaderType)
}

func (t *Tool) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf("Read data from %s, the input must be json schema: %s", t.ReaderType, string(bytes)) + `
Example Input: {"url": "https://www.example.com/pdf"}`
}

func (t *Tool) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"url": {
				Type:        tool.TypeString,
				Description: "The url to read",
			},
		},
		Required: []string{"url"},
	}
}

func (t *Tool) Strict() bool {
	return true
}

func (t *Tool) Call(_ context.Context, input string) (string, error) {
	var param ReadParam

	err := json.Unmarshal([]byte(input), &param)
	if err != nil {
		return "json unmarshal error, please try again", nil
	}

	text, err := t.Reader.Read(param.Url)
	if err != nil {
		return err.Error(), nil
	}

	return text, nil
}
