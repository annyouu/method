package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// for文内のdefer文を検出する
var forDeferAnalyzer = &analysis.Analyzer{
	Name: "fordefer",
	Doc: "forとdeferを検知します",
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// 各ファイルのASTを走査
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// for文の場合
			if forStmt, ok := n.(*ast.ForStmt); ok {
				// for文下のステートメントをチェック
				for _, stmt := range forStmt.Body.List {
					if deferStmt, ok := stmt.(*ast.DeferStmt); ok {
						pass.Reportf(deferStmt.Pos(), "defer文が検出されました")
					}
				}
				return true
			}

			// for range文の場合
			if rangeStmt, ok := n.(*ast.RangeStmt); ok {
				for _, stmt := range rangeStmt.Body.List {
					if deferStmt, ok := stmt.(*ast.DeferStmt); ok {
						pass.Reportf(deferStmt.Pos(), "defer文が検出されました")
					}
				}
				return true
			}
			return true
		})
	}
	return nil, nil
}

func main() {
	singlechecker.Main(forDeferAnalyzer)
}