package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

type SingleFileVisitor struct {
	file *FileVisitor
}

func (s *SingleFileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	file, ok := node.(*ast.File)
	if !ok {
		return s
	}
	s.file = &FileVisitor{
		Package: file.Name.String(),
	}
	return s.file
}

func (s *SingleFileVisitor) Get() File {
	types := make([]Type, 0, len(s.file.Types))
	for _, t := range s.file.Types {
		types = append(types, Type{
			Name:   t.name,
			Fields: t.fields,
		})
	}
	return File{
		Package: s.file.Package,
		Imports: s.file.Imports,
		Types:   types,
	}
}

type FileVisitor struct {
	Package string
	Imports []string
	Types   []*TypeVisitor
}

func (f *FileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.GenDecl:
		if n.Tok == token.IMPORT {
			for _, spec := range n.Specs {
				f.Imports = append(f.Imports, spec.(*ast.ImportSpec).Path.Value)
			}
		}
	case *ast.TypeSpec:
		v := &TypeVisitor{name: n.Name.String()}
		f.Types = append(f.Types, v)
		return v
	}
	return f
}

type File struct {
	Package string
	Imports []string
	Types   []Type
}

type TypeVisitor struct {
	name   string
	fields []Field
}

func (t *TypeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	n, ok := node.(*ast.Field)
	if !ok {
		return t
	}
	var fieldTyp string
	switch typ := n.Type.(type) {
	case *ast.Ident:
		fieldTyp = typ.String()
	case *ast.IndexExpr:
		fieldTyp = typ.X.(*ast.Ident).Name + "[" + typ.Index.(*ast.Ident).Name + "]"
	case *ast.SelectorExpr:
		fieldTyp = typ.X.(*ast.Ident).Name + "." + typ.Sel.Name
	case *ast.StarExpr:
		fieldTyp = typ.X.(*ast.SelectorExpr).X.(*ast.Ident).Name + "." + typ.X.(*ast.SelectorExpr).Sel.Name
	case *ast.ArrayType:
		fieldTyp = "[]byte"
	default:
		panic(fmt.Sprintf("无法识别的类型: %v\n", typ))
	}

	// 是node.Names 不是 node.Name的原因： 存在这种写法： a,b,c string
	for _, name := range n.Names {
		t.fields = append(t.fields, Field{
			Name: name.String(),
			Type: fieldTyp,
		})
	}

	return t
}

type Type struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}
