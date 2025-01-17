package wikipedia

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/antgroup/aievo/tool"
	"github.com/antgroup/aievo/utils/json"
)

var ErrUnexpectedAPIResult = errors.New("unexpected result from wikipedia api")

type Tool struct {
	// The number of wikipedia pages to include in the result.
	TopK int
	// The number of characters to take from each page.
	DocMaxChars int
	// The language code to use.
	LanguageCode string
	// The user agent sent in the header. See https://www.mediawiki.org/wiki/API:Etiquette.
	UserAgent string
}

var _ tool.Tool = &Tool{}

// New creates a new wikipedia tool to find wikipedia pages using the wikipedia api. TopK is set
// to 2, DocMaxChars is set to 2000 and the language code is set to "en".
func New(opts ...Option) (*Tool, error) {
	options := NewDefaultOptions()

	for _, opt := range opts {
		opt(options)
	}

	return &Tool{
		TopK:         options.TopK,
		DocMaxChars:  options.DocMaxChars,
		LanguageCode: options.LanguageCode,
		UserAgent:    options.UserAgent,
	}, nil
}

func (t *Tool) Name() string {
	return "Wikipedia"
}

func (t *Tool) Description() string {
	return `A wrapper around Wikipedia. 
	Useful for when you need to answer general questions about 
	people, places, companies, facts, historical events, or other subjects. 
	Input should be a search query.`
}

func (t *Tool) Schema() *tool.PropertiesSchema {
	return &tool.PropertiesSchema{
		Type: tool.TypeJson,
		Properties: map[string]tool.PropertySchema{
			"query": {
				Type:        tool.TypeString,
				Description: "The query to search for wikipedia pages.",
			},
		},
		Required: []string{"query"},
	}
}

func (t *Tool) Strict() bool {
	return true
}

func (t *Tool) Call(ctx context.Context, input string) (string, error) {
	var m map[string]interface{}

	err := json.Unmarshal([]byte(input), &m)
	if err != nil {
		return "json unmarshal error, please try agent", nil
	}

	if m["query"] == "" {
		return "query is required", nil
	}

	result, err := t.searchWiki(ctx, m["query"].(string))
	if err != nil {
		fmt.Println(err)
		return "Query Wikipedia Error, Please try again!", nil
	}
	return result, nil
}

func (t *Tool) searchWiki(ctx context.Context, input string) (string, error) {
	searchResult, err := search(ctx, t.TopK, input, t.LanguageCode, t.UserAgent)
	if err != nil {
		return "", err
	}

	if len(searchResult.Query.Search) == 0 {
		return "no wikipedia pages found", nil
	}

	result := ""

	for _, search := range searchResult.Query.Search {
		getPageResult, err := getPage(ctx, search.PageID, t.LanguageCode, t.UserAgent)
		if err != nil {
			return "", err
		}

		page, ok := getPageResult.Query.Pages[strconv.Itoa(search.PageID)]
		if !ok {
			return "", ErrUnexpectedAPIResult
		}
		if len(page.Extract) >= t.DocMaxChars {
			result += page.Extract[0:t.DocMaxChars]
			continue
		}
		result += page.Extract
	}

	return result, nil
}
