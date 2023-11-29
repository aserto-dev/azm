package parser_test

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/parser"
)

type TreeShapeListener struct {
	*parser.BaseAzmListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func (this *TreeShapeListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func TestParser(t *testing.T) {
	input, _ := antlr.NewFileStream("./parser_test.txt")
	lexer := parser.NewAzmLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.Prog()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}
