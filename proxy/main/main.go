package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"rpc-go-netty/proxy"
	"rpc-go-netty/proxy/service"
)

func main() {

	clientProxy := proxy.NewRemoteClientProxy(service.MyService())
	clientProxy.Invoke("Hello", nil)

}

type ServiceRegistrar struct {
	services []interface{}
}

func (r *ServiceRegistrar) Register(service interface{}) {
	// 在这里执行注册逻辑
	r.services = append(r.services, service)
}

func (r *ServiceRegistrar) ScanAndRegister(packageName string) {
	pkg, err := build.Import(packageName, "", build.FindOnly)
	if err != nil {
		fmt.Println("Failed to import package:", err)
		return
	}

	fmt.Println(pkg.GoFiles)

	for _, filename := range pkg.GoFiles {
		fmt.Println(filename)
		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, filepath.Join(pkg.Dir, filename), nil, parser.AllErrors)
		if err != nil {
			fmt.Println("Failed to parse file:", err)
			continue
		}

		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := typeSpec.Type.(*ast.StructType); ok {
							// 检查结构体类型是否带有 ServiceAnnotation 注解
							if hasServiceAnnotation(typeSpec.Type) {
								// 创建结构体类型的实例，并注册服务
								serviceType := reflect.TypeOf(typeSpec.Name.Name)
								fmt.Println(serviceType)
								service := reflect.New(serviceType).Elem().Interface()
								fmt.Println(service)
								r.Register(service)
							}
						}
					}
				}
			}
		}
	}
}

// 检查结构体类型是否带有 ServiceAnnotation 注解
func hasServiceAnnotation(typ ast.Expr) bool {
	var (
		annotationTypeName = "ServiceAnnotation" // 注解类型名称
		found              = false
	)

	ast.Inspect(typ, func(node ast.Node) bool {
		if ident, ok := node.(*ast.Ident); ok && ident.Name == annotationTypeName {
			found = true
			return false
		}
		return true
	})

	return found
}
