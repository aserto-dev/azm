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

func (l *TreeShapeListener) ExitUnionRel(c *parser.UnionRelContext) {
	fmt.Println("ExitUnionRel", c.GetText())
}

func (l *TreeShapeListener) ExitIntersectRel(c *parser.IntersectRelContext) {
	fmt.Println("ExitIntersectRel", c.GetText())
}

func (l *TreeShapeListener) ExitExclusionRel(c *parser.ExclusionRelContext) {
	fmt.Println("ExitExclusionRel", c.GetText())
}

func (l *TreeShapeListener) ExitSingleRel(c *parser.SingleRelContext) {
	fmt.Println("ExitSingleRel", c.GetText())
}

func (l *TreeShapeListener) ExitWildcardRel(c *parser.WildcardRelContext) {
	fmt.Println("ExitWildcardRel", c.GetText())
}

func (l *TreeShapeListener) ExitSubjectRel(c *parser.SubjectRelContext) {
	fmt.Println("ExitSubjectRel", c.GetText())
}

func (l *TreeShapeListener) ExitArrowRel(c *parser.ArrowRelContext) {
	fmt.Println("ExitArrowRel", c.GetText())
}

func TestParser(t *testing.T) {
	// input := antlr.NewInputStream("user | user:* | group#member | parent->viewer\n")
	input, _ := antlr.NewFileStream("./parser_test.txt")
	lexer := parser.NewAzmLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.Prog()
	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}
