package arxiv

import (
	"context"
	"errors"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

type Tool struct {
	// client is the arXiv client.
	client *Client
}

var _ tool.Tool = &Tool{}

// New creates a new Arxiv Tool.
func New(opts ...Option) (*Tool, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &Tool{
		client: NewClient(
			options.Topk,
			options.UserAgent,
		),
	}, nil
}

// Name returns the name of the tool.
func (t *Tool) Name() string {
	return "Arxiv Search"
}

// Description returns the description of the tool.
func (t *Tool) Description() string {
	bytes, _ := json.Marshal(t.Schema())
	return `A wrapper around arXiv Search API.
Search for scientific papers on arXiv, the input must be json schema:` + string(bytes) + `
Example Input: {\"query\": \"machine learning, LLM, AI\"}`
}

func (t *Tool) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"query": {
				Type:        tool.TypeString,
				Description: "The query to search for papers on arXiv, must be English",
			},
		},
		Required: []string{"query"},
	}
}

func (t *Tool) Strict() bool {
	return true
}

// Call searches for papers on arXiv.
// Input should be a search query.
func (t *Tool) Call(ctx context.Context, input string) (string, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(input), &m)
	if err != nil {
		return "json unmarshal error, please try agent", nil
	}
	if m["query"] == nil {
		return "query is required", nil
	}
	result, err := t.client.Search(ctx, m["query"].(string))
	if err != nil {
		if errors.Is(err, ErrNoGoodResult) {
			return "No good arXiv Search Results were found", nil
		}
		return "", err
	}
	return result, nil
}
