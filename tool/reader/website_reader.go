package reader

import (
	"time"

	"github.com/go-shiori/go-readability"
)

type WebsiteReader struct {
}

var _ Reader = &WebsiteReader{}

func NewWebsiteReader() Reader {
	return &WebsiteReader{}
}

func (r *WebsiteReader) Read(url string) (string, error) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return "", err
	}
	return article.TextContent, nil
}
