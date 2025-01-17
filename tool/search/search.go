package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

const (
	_defaultTopK   = 6
	_defaultEngine = "google"
)

var searchEngines = map[string]Factory{
	"google": NewGoogleSearch,
	"bing":   NewBingSearch,
}

type Tool struct {
	TopK   int
	Engine string
	client *Client
}

var _ tool.Tool = &Tool{}

func New(opts ...Option) (*Tool, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	factory, ok := searchEngines[options.Engine]
	if !ok {
		factory = searchEngines[_defaultEngine]
	}
	if options.TopK <= 0 {
		options.TopK = _defaultTopK
	}
	return &Tool{
		TopK:   options.TopK,
		Engine: options.Engine,
		client: factory(options.ApiKey),
	}, nil
}

func (t *Tool) Name() string {
	return fmt.Sprintf("%s Search", strings.ToUpper(t.Engine))
}

func (t *Tool) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return fmt.Sprintf(`A wrapper around %s Search.
Useful for when you need to answer questions about current events, 
the input must be json schema: %s`, strings.ToUpper(t.Engine), string(bytes)) + `
Example Input: {\"query\": \"machine learning, LLM, AI\"}`
}

func (t *Tool) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"query": {
				Type:        tool.TypeString,
				Description: fmt.Sprintf("the query to search on %s, must be English", strings.ToUpper(t.Engine)),
			},
		},
		Required: []string{"query"},
	}
}

func (t *Tool) Strict() bool {
	return true
}

func (t *Tool) Call(_ context.Context, input string) (string, error) {
	var m map[string]interface{}

	err := json.Unmarshal([]byte(input), &m)
	if err != nil {
		return "json unmarshal error, please try agent", nil
	}

	if m["query"] == "" {
		return "query is required", nil
	}

	ret, err := t.client.Search(m["query"].(string), t.TopK)
	if err != nil {
		return "Query Search Engine Error, Please Try Again", nil
	}
	return ret, nil
}
