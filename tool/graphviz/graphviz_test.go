package graphviz

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/antgroup/aievo/llm/openai"
)

func TestGraphviz(t *testing.T) {
	client, err := openai.New(
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		openai.WithModel(os.Getenv("OPENAI_MODEL")),
		openai.WithBaseURL(os.Getenv("OPENAI_BASE_URL")))
	if err != nil {
		log.Fatal(err)
	}

	tool, err := NewGraphvizTool(client, 3)
	if err != nil {
		panic(err)
	}
	output, err := tool.Call(context.Background(), `
{
	"sop": "1. 产品根据用户的需求，写相关的需求文档，然后交给架构师
2. 架构师根据需求文档，写出系统设计的文档，然后交给项目经理
3. 项目经理根据系统设计文档，进行任务划分，然后分发任务给相应的程序员
4. 程序员A/B/C根据任务文档，写代码，交给测试员进行测试
5. 测试员对代码进行测试，如果通过，则结束流程，如果不通过，打回给程序员修改，修改后还需要交给测试员测试，重复该流程
6. 模块A/B/C都通过测试，项目完成",
	"agent_descriptions": "产品: 产品agent
架构师: 架构师agent
程序员: 程序员agent
测试: 测试agent"
}
`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("output: ")
	fmt.Println(output)
}
