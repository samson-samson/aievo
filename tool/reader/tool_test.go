package reader

import (
	"context"
	"fmt"
	"testing"

	"github.com/antgroup/aievo/utils/json"
)

func TestWebsiteReader(t *testing.T) {
	tool, _ := NewReader(
		WithReaderType("website"),
	)
	var param ReadParam
	param.Url = "https://53ai.com/news/qianyanjishu/901.html"
	bytes, _ := json.Marshal(param)
	text, _ := tool.Call(context.Background(), string(bytes))
	fmt.Println(text)
}

func TestPdfReader(t *testing.T) {
	tool, _ := NewReader(
		WithReaderType("pdf"),
	)
	var param ReadParam
	param.Url = "https://arxiv.org/pdf/2308.00352"
	bytes, _ := json.Marshal(param)
	text, _ := tool.Call(context.Background(), string(bytes))
	fmt.Println(text)
}
