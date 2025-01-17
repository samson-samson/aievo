package code

import (
	"context"
	"fmt"
	"testing"
)

func TestRunGolangCode(t *testing.T) {
	code := `
	package main

	import (
		"fmt"
	)

	func main() {
		fmt.Println("hello world")
	}
	`
	tool, _ := New(WithProgramLangType("golang"))
	result, _ := tool.Call(context.Background(), fmt.Sprintf(`{
	"code": "%s"
}`, code))
	fmt.Println(result)
}

func TestRunPythonCode(t *testing.T) {
	code := `
print("hello world")
	`
	tool, _ := New(WithProgramLangType("python"))
	result, _ := tool.Call(context.Background(), fmt.Sprintf(`{
	"code": "%s"
}`, code))
	fmt.Println(result)
}

func TestRunJavaCode(t *testing.T) {
	code := `
public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello World!");
    }
}
	`
	tool, _ := New(WithProgramLangType("java"))
	result, _ := tool.Call(context.Background(), fmt.Sprintf(`{
	"code": "%s"
}`, code))
	fmt.Println(result)
}
