package reader

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/antgroup/aievo/utils/pdf"
)

type PdfReader struct {
}

var _ Reader = &PdfReader{}

func NewPdfReader() Reader {
	return &PdfReader{}
}

func (r *PdfReader) Read(url string) (string, error) {
	// check poppler version
	err := pdf.CheckPopplerVersion()
	if err != nil {
		return err.Error(), nil
	}

	// download pdf file
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	// extract text from pdf file
	bytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}
	pages, err := pdf.ExtractOrError(bytes)
	if err != nil {
		return "", err
	}
	texts := make([]string, 0)
	for _, page := range pages {
		texts = append(texts, page.Content)
	}

	return strings.Join(texts, "\n"), nil
}
