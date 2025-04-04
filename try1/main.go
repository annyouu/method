package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)


func countVarsAndConsts(f *ast.File) (varCount, constCount int) {
	for _, decl := range f.Decls {
		// genDeclはインポート、変数、定数、型の宣言を表す
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch genDecl.Tok {
		case token.VAR:
			for _, spec := range genDecl.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					varCount += len(vs.Names)
				}
			}
		case token.CONST:
			for _, spec := range genDecl.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					constCount += len(vs.Names)
				}
			}
		}
	}
	return
}


func main() {
	src := `
		package main

		import "fmt"

		var msg string
		var a, b, c int
		const Pi = 3.14
		const x, y = 1, 2

		func main() {
			msg = "Hello"
			name := "Gopher"
			fmt.Println(msg, name)
		}
	`

	// 位置情報を管理する
	fset := token.NewFileSet()

	// Goファイルをパッケージ単位で解析
	node, err := parser.ParseFile(fset, "example.go", src, parser.AllErrors)
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	// 変数、定数を数える
	varCount, constCount := countVarsAndConsts(node)

	// 結果を表示
	fmt.Println("=== 解析結果 ===")
	fmt.Println("ファイル: example.go")
	fmt.Println("変数の数:", varCount)
	fmt.Println("定数の数:", constCount)
}