package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// ファイル内の変数と定数の数を数える
func constDeclsInFile(f *ast.File) (varCount, constCount int) {
	for _, decl := range f.Decls {
		// genDeclはインポート、変数、定数、型の宣言を表す
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch genDecl.Tok {
		case token.VAR:
			// ValueSpecのNamesフィールドに宣言された識別子が格納されている
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
	// 解析対象のディレクトリパスを指定(絶対パスで指定する)
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("カレントディレクトリ取得エラー:", err)
		os.Exit(1)
	}

	dirPath := filepath.Join(wd, "example")

	// token.FileSetを作成(位置情報管理)
	fset := token.NewFileSet()

	// 指定ディレクトリ内のGoファイルをパッケージ単位で解析
	pkgs, err := parser.ParseDir(fset, dirPath, nil, 0)
	if err != nil {
		fmt.Println("エラー:", err)
		os.Exit(1)
	}

	// 各パッケージごとに処理する
	for pkgName, pkg := range pkgs {
		fmt.Printf("===== パッケージ: %s =====\n", pkgName)
		// 各ファイルごとに変数・定数の数を数える
		for fileName, file := range pkg.Files {
			varCount, constCount := constDeclsInFile(file)
			absPath, _ := filepath.Abs(fileName)
			fmt.Printf("ファイル: %s\n", absPath)
			fmt.Printf("変数の数: %d\n", varCount)
			fmt.Printf("定数の数: %d\n", constCount)
		}
	}
}