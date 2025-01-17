package arxiv

import (
	"context"
	"fmt"
	"testing"
)

func TestArxiv(t *testing.T) {
	tool, _ := New(
		WithTopk(10),
		WithUserAgent("aievo/1.0"),
	)
	result, _ := tool.Call(context.Background(), `{
	"query": "ai agent"
}`)
	fmt.Println(result)
}
