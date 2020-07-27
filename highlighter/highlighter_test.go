package highlighter_test

import (
	"github.com/gookit/highlighter"
	"testing"
)

var src = []byte(`package main
import "fmt"
func main() {
	fmt.Println("Hey there, Go.")
}
`)

func TestSourceCode(t *testing.T) {
	highlighter.SourceCode(string(src))
}
