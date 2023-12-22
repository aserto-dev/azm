package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/types"
)

type RelationRef = types.RelationRef
type Permission = types.Permission
type PermissionTerm = types.PermissionTerm
type RelationName = types.RelationName
type ExclusionPermission = types.ExclusionPermission
type ObjectName = types.ObjectName

func ParseRelation(input string) []*RelationRef {
	p := newParser(input)
	rTree := p.Relation()

	var v RelationVisitor
	return v.Visit(rTree).([]*RelationRef)
}

func ParsePermission(input string) *Permission {
	p := newParser(input)
	pTree := p.Permission()

	var v PermissionVisitor
	return v.Visit(pTree).(*Permission)
}

func newParser(input string) *AzmParser {
	lexer := NewAzmLexer(antlr.NewInputStream(input))
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	return p
}
