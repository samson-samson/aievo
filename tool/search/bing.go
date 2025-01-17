package search

import (
	"encoding/json"
	"errors"
	"strconv"
)

var ErrorEncodeBingResp = errors.New("error encoding Bing search result")

// NewBingSearch creates search for Bing
func NewBingSearch(apiKey string) *Client {
	return NewSearch("bing",
		apiKey,
		&BingSearchParamsHandler{},
		&BingSearchResultHandler{
			RequiredField: "organic_results",
		})
}

type BingSearchParamsHandler struct {
}

func (h *BingSearchParamsHandler) Handle(input string, pageIndex, pageSize int) (map[string]string, error) {
	return map[string]string{
		"q":     input,
		"first": strconv.Itoa((pageIndex-1)*pageSize + 1),
		"count": strconv.Itoa(pageSize),
	}, nil
}

type BingSearchResultHandler struct {
	RequiredField string
}

func (h *BingSearchResultHandler) GetRequiredField() string {
	return h.RequiredField
}

func (h *BingSearchResultHandler) Handle(result string) (string, error) {
	var BingDatas []BingData
	err := json.Unmarshal([]byte(result), &BingDatas)
	if err != nil {
		return "", ErrorEncodeBingResp
	}
	var output string
	for _, data := range BingDatas {
		output += Title + data.Title + "\n"
		output += Snippet + data.Snippet + "\n"
		output += Link + data.Link + "\n\n"
	}
	return output, nil
}
