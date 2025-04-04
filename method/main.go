// パッケージ関数の一覧を表示する方法
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	src := `
		package main

		import "fmt"

		func init() {
			fmt.Println("init 1")
		}

		func init() {
			fmt.Println("init 2")
		}

		func _() {
			fmt.Println("black function")
		}
		
		func Hello(name string) string {
			return "Hello, " + name
		}

		func main() {
			Hello("World")
		}
	`

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "example.go", src, parser.AllErrors)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	// ASTを走査してパッケージ関数(Recv == nil)の名前を表示
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			// Recvがnilならパッケージ関数
			if fn.Recv == nil {
				fmt.Printf("関数名: %s\n", fn.Name.Name)
			}
		}
	}
}