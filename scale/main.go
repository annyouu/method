package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

// 対象のGoコードを埋め込む
const src = `
package main

import "fmt"

// 短い関数
func shortFunc() {
	fmt.Println("short")
}

// 中くらいの関数
func mediumFunc() {
	fmt.Println("line1")
	fmt.Println("line2")
	fmt.Println("line3")
}

// 長い関数
func longFunc() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
	fmt.Println("done")
}
`

// 関数情報をまとめる構造体
type FuncInfo struct {
	Name string `json:"name"`
	File string `json:"file"`
	Line int `json:"line"`
	LineSize int `json:"line_size"`
	ByteSize int `json:"byte_size"`
	StmtSize int `json:"stmt_size"`
}

func main() {
	// ソースコードをパース
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, "src.go", src, parser.AllErrors)
	if err != nil {
		fmt.Println("構文エラー:", err)
		return
	}

	var funcs []FuncInfo

	// 関数を収集
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		// メソッドは除外
		if funcDecl.Recv != nil {
			continue
		}

		start := fs.Position(funcDecl.Pos())
		end := fs.Position(funcDecl.End())
		lineCount := end.Line - start.Line + 1
		byteSize := funcDecl.End() - funcDecl.Pos()
		stmtCount := len(funcDecl.Body.List)

		funcs = append(funcs, FuncInfo{
			Name: funcDecl.Name.Name,
			File: start.Filename,
			Line: start.Line,
			LineSize: lineCount,
			ByteSize: int(byteSize),
			StmtSize: stmtCount,
		})
	}

	// 行数でソート(大きい順)
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].LineSize > funcs[j].LineSize
	})

	// JSONに変換して表示
	jsonBytes, err := json.MarshalIndent(funcs, "", " ")
	if err != nil {
		fmt.Println("JSON変換エラー:", err)
		return
	}

	fmt.Println("JSON形式の出力:")
	fmt.Println(string(jsonBytes))

	fmt.Println("関数のサイズ(大きい順):")
	for _, f := range funcs {
		fmt.Printf("関数名: %-12s 行数: %3d バイト数: %4d 文数: %2d (%s:%d)\n", 
		f.Name, f.LineSize, f.ByteSize, f.StmtSize, f.File, f.Line)
	}
}