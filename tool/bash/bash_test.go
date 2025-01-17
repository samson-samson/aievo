package bash

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestBash(t *testing.T) {
	shellTool, _ := New()
	output, err := shellTool.Call(context.Background(), `{
	"command": "echo hello world"
}`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output:")
	fmt.Println(output)
}
