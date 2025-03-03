package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed tpl.gohtml
var genOrm string

func gen(w io.Writer, srcFile string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	s := &SingleFileVisitor{}
	ast.Walk(s, f)
	file := s.Get()
	tpl := template.New("gen-orm")
	tpl, err = tpl.Parse(genOrm)
	if err != nil {
		return err
	}
	return tpl.Execute(w, &Data{
		File: file,
		Opts: []string{"Lt", "Eq", "Gt"},
	})
}

type Data struct {
	File
	Opts []string
}

/*
命令行跑， 生成testdata/user_gen.go

1. cd orm_gen // 注意：package 需要改成main
2. go install . // go install 到本地，这行代码会在gopath/bin下生成可执行命令，通过文件夹名称调用
3. cd testdata // testdata 默认会被go忽略
4. orm_gen user.go // user.go 是参数
注意：main.go中的import "html/template" 需要改成 "text/template"， 不然“”会被转义
*/
func main() {
	src := os.Args[1]
	dstDir := filepath.Dir(src)
	fileName := filepath.Base(src)
	idx := strings.LastIndexByte(fileName, '.')
	dst := filepath.Join(dstDir, fileName[:idx]+"_gen.go")
	f, err := os.Create(dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	err = gen(f, src)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("生成成功")
}
