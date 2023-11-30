package parser_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/parser"
)

type production int

const (
	unknown production = iota
	permission
	relation
)

var (
	errNotPermissionProduction = errors.New("not a permission production")
	errNotRelationProduction   = errors.New("not a relation production")
	errProductionAlreadySet    = errors.New("production already set")
	errProductionNotSet        = errors.New("production not set (unknown)")
)

type AzmListener struct {
	*parser.BaseAzmListener
	production production
	relations  []*model.ObjectRelation
}

func NewAzmListener() *AzmListener {
	azm := new(AzmListener)
	azm.production = unknown
	azm.relations = []*model.ObjectRelation{}
	return azm
}

func (l *AzmListener) GetPermission() ([]*model.ObjectRelation, error) {
	if l.production == permission {
		return l.relations, nil
	}
	return nil, errNotPermissionProduction
}

func (l *AzmListener) GetRelation() ([]*model.ObjectRelation, error) {
	if l.production == relation {
		return l.relations, nil
	}
	return nil, errNotRelationProduction
}

func (l *AzmListener) EnterPermission(c *parser.PermissionContext) {
	if l.production == unknown {
		l.production = permission
		return
	}
	panic(errProductionAlreadySet)
}

func (l *AzmListener) EnterRelation(c *parser.RelationContext) {
	if l.production == unknown {
		l.production = relation
		return
	}
	panic(errProductionAlreadySet)
}

func (l *AzmListener) ExitUnion(c *parser.UnionContext) {
	// relation production
	if l.production == relation {
		fmt.Println("ExitUnionRel", c.GetText())
		return
	}

	// permission production
	if l.production == permission {
		fmt.Println("ExitUnionRel", c.GetText())
		return
	}
}

func (l *AzmListener) ExitIntersection(c *parser.IntersectionContext) {
	// permission production
	if l.production == permission {
		fmt.Println("ExitIntersectRel", c.GetText())
	}
}

func (l *AzmListener) ExitExclusion(c *parser.ExclusionContext) {
	// permission production
	if l.production == permission {
		fmt.Println("ExitExclusionRel", c.GetText())
		return
	}
}

func (l *AzmListener) ExitSingleRel(c *parser.SingleRelContext) {
	fmt.Println("ExitSingleRel", c.GetText())
	l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
}

func (l *AzmListener) ExitWildcardRel(c *parser.WildcardRelContext) {
	fmt.Println("ExitWildcardRel", c.GetText())
	l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
}

func (l *AzmListener) ExitSubjectRel(c *parser.SubjectRelContext) {
	fmt.Println("ExitSubjectRel", c.GetText())
	l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
}

func (l *AzmListener) ExitArrowRel(c *parser.ArrowRelContext) {
	fmt.Println("ExitArrowRel", c.GetText())
	l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
}

func (l *AzmListener) ExitRelation(c *parser.RelationContext) {
	fmt.Println("ExitRelation", c.GetText())
}

func (l *AzmListener) ExitPermission(c *parser.PermissionContext) {
	fmt.Println("ExitPermission", c.GetText())
}

func TestParser(t *testing.T) {
	input := antlr.NewInputStream("user | user:* | group#member | parent->viewer")
	// input, _ := antlr.NewFileStream("./parser_test.txt")
	lexer := parser.NewAzmLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	listener := NewAzmListener()

	rTree := p.Relation()
	antlr.ParseTreeWalkerDefault.Walk(listener, rTree)
	rel, err := listener.GetRelation()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", rel)

	// pTree := p.Permission()
	// antlr.ParseTreeWalkerDefault.Walk(listener, pTree)
}
