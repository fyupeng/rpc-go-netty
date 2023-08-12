package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

// 注解类型
type Annotation struct {
	Name      string
	Arguments []string
}

// 自动注册函数
func AutoRegister(packagePath string) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)

				if !ok {
					return true
				}

				if !ok || fn.Doc == nil || len(fn.Doc.List) == 0 {
					return true
				}

				// 如果找到了函数声明
				fmt.Println("Function Name:", fn.Name.Name)
				if fn.Type.Params != nil {
					// 打印函数参数列表
					fmt.Println("Function Parameters:")
					for _, param := range fn.Type.Params.List {
						names := make([]string, len(param.Names))
						for i, ident := range param.Names {
							names[i] = ident.Name
						}
						fmt.Printf("  %s: %s\n", strings.Join(names, ", "), param.Type)
					}
				}

				if fn.Type.Results != nil {
					// 打印函数返回值
					fmt.Println("Function Results:")
					for _, result := range fn.Type.Results.List {
						fmt.Println("  ", result.Type)
					}
				}

				for _, comment := range fn.Doc.List {
					text := comment.Text[3:] // 去除注释前面的 "//"
					annotation := parseAnnotation(text)
					if annotation != nil {
						fmt.Printf("Found annotation: %s\n", annotation.Name)
						fmt.Printf("Found annotation: %s\n", annotation.Arguments)
						// 在这里执行自动注册的逻辑，可以根据需要添加你的业务逻辑
					}
				}

				return true
			})
		}
	}
}

// 解析注解
func parseAnnotation(comment string) *Annotation {
	// 在这里添加解析注解的逻辑，可以根据你的注解格式进行解析
	// 这里只是一个示例，假设注解格式为 @Register("cn.fyupeng.service", "arg1", "arg2")
	if comment[:9] != "@Register" {
		return nil
	}

	args := extractArguments(comment[9:])
	annotation := &Annotation{
		Name:      "Register",
		Arguments: args,
	}

	return annotation
}

// 提取注解参数
func extractArguments(comment string) []string {
	// 假设参数使用英文逗号分隔
	args := make([]string, 0)
	start := -1
	for i := 0; i < len(comment); i++ {
		if comment[i] == '"' && (i == 0 || comment[i-1] != '\\') {
			if start == -1 {
				start = i + 1
			} else {
				args = append(args, comment[start:i])
				start = -1
			}
		}
	}

	return args
}

func main() {
	AutoRegister(os.Getenv("GOPATH") + "/rpc-go-netty/test/demo")
}
