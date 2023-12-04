package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
)

func ParseRelation(input string) []*model.RelationTerm {
	p := newParser(input)
	rTree := p.Relation()

	var v RelationVisitor
	return v.Visit(rTree).([]*model.RelationTerm)
}

func ParsePermission(input string) *model.Permission {
	p := newParser(input)
	pTree := p.Permission()

	var v PermissionVisitor
	return v.Visit(pTree).(*model.Permission)
}

func newParser(input string) *AzmParser {
	lexer := NewAzmLexer(antlr.NewInputStream(input))
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	return p
}
