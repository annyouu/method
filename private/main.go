package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)


// テスト用のコードをパースして静的解析する
const src = `
package sample

func ExportedFunc() {}
func internalFunc() {}
func anotherInternal() {}

type MyType struct{}

func (m MyType) MethodFunc() {}
`

func main() {
	// ソースコードのパース
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, "sample.go", src, 0)
	if err != nil {
		fmt.Println("パースエラー:", err)
		return
	}

	// ASTを走査して、非公開な関数を抽出
	fmt.Println("非公開なパッケージ関数一覧:")
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// メソッドを除外
		if funcDecl.Recv != nil {
			continue
		}

		// 非公開な関数
		if !funcDecl.Name.IsExported() {
			pos := fs.Position(funcDecl.Pos())
			fmt.Printf("%s (at %s)\n", funcDecl.Name.Name, pos)
		}
	}
}
