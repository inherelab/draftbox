// Package highlighter is a command line code highlighting for go codes.
package highlighter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

const (
	STRING  = "string"
	COMMENT = "comment"
	DEFAULT = "default"
	TYPE    = "type"
	KEYWORD = "keyword"
)

var defTheme = map[string]string{
	STRING:  "cyan",
	COMMENT: "gray",
	TYPE:    "yellow",
	KEYWORD: "green",
}

// Highlighter definition
type Highlighter struct {
}

func SourceCode(code string) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", code, 0)
	if err != nil {
		panic(err)
	}

	// Create an ast.CommentMap from the ast.File's comments.
	// This helps keeping the association between comments
	// and AST nodes.
	cm := ast.NewCommentMap(fset, f, f.Comments)

	// Remove the first variable declaration from the list of declarations.
	for i, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)

		if ok {
			fmt.Println(gen.Tok.String())
		}

		if ok && gen.Tok == token.VAR {
			copy(f.Decls[i:], f.Decls[i+1:])
			f.Decls = f.Decls[:len(f.Decls)-1]
		}
	}

	// Use the comment map to filter comments that don't belong anymore
	// (the comments associated with the variable declaration), and create
	// the new comments list.
	f.Comments = cm.Filter(f).Comments()

	// Print the modified AST.
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}

	fmt.Printf("%s", buf.Bytes())
}

func SourceBytes(bs []byte) {

}

func SourceFile(file string) {

}
