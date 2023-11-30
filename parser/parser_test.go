package parser_test

import (
	"errors"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/parser"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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

// type AzmListener struct {
//     *parser.BaseAzmListener
//     production production
//     relations  []*model.ObjectRelation
// }

// func NewAzmListener() *AzmListener {
//     azm := new(AzmListener)
//     azm.production = unknown
//     azm.relations = []*model.ObjectRelation{}
//     return azm
// }

// func (l *AzmListener) GetPermission() ([]*model.ObjectRelation, error) {
//     if l.production == permission {
//         return l.relations, nil
//     }
//     return nil, errNotPermissionProduction
// }

// func (l *AzmListener) GetRelation() ([]*model.ObjectRelation, error) {
//     if l.production == relation {
//         return l.relations, nil
//     }
//     return nil, errNotRelationProduction
// }

// func (l *AzmListener) EnterPermission(c *parser.PermissionContext) {
//     if l.production == unknown {
//         l.production = permission
//         return
//     }
//     panic(errProductionAlreadySet)
// }

// func (l *AzmListener) EnterRelation(c *parser.RelationContext) {
//     if l.production == unknown {
//         l.production = relation
//         return
//     }
//     panic(errProductionAlreadySet)
// }

// func (l *AzmListener) ExitUnion(c *parser.UnionContext) {
//     // relation production
//     if l.production == relation {
//         fmt.Println("ExitUnionRel", c.GetText())
//         return
//     }

//     // permission production
//     if l.production == permission {
//         fmt.Println("ExitUnionRel", c.GetText())
//         return
//     }
// }

// func (l *AzmListener) ExitIntersection(c *parser.IntersectionContext) {
//     // permission production
//     if l.production == permission {
//         fmt.Println("ExitIntersectRel", c.GetText())
//     }
// }

// func (l *AzmListener) ExitExclusion(c *parser.ExclusionContext) {
//     // permission production
//     if l.production == permission {
//         fmt.Println("ExitExclusionRel", c.GetText())
//         return
//     }
// }

// func (l *AzmListener) ExitSingleRel(c *parser.SingleRelContext) {
//     fmt.Println("ExitSingleRel", c.GetText())
//     l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
// }

// func (l *AzmListener) ExitWildcardRel(c *parser.WildcardRelContext) {
//     fmt.Println("ExitWildcardRel", c.GetText())
//     l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
// }

// func (l *AzmListener) ExitSubjectRel(c *parser.SubjectRelContext) {
//     fmt.Println("ExitSubjectRel", c.GetText())
//     l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
// }

// func (l *AzmListener) ExitArrowRel(c *parser.ArrowRelContext) {
//     fmt.Println("ExitArrowRel", c.GetText())
//     l.relations = append(l.relations, &model.ObjectRelation{Object: "", Relation: ""})
// }

// func (l *AzmListener) ExitRelation(c *parser.RelationContext) {
//     fmt.Println("ExitRelation", c.GetText())
// }

// func (l *AzmListener) ExitPermission(c *parser.PermissionContext) {
//     fmt.Println("ExitPermission", c.GetText())
// }

type RelationVisitor struct {
	parser.BaseAzmVisitor
}

func (v *RelationVisitor) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case *parser.PermissionContext:
		panic("RelationVisitor cannot visit permissions")
	case *parser.RelationContext:
		return t.Accept(v)
	}

	return nil
}

func (v *RelationVisitor) VisitRelation(c *parser.RelationContext) interface{} {
	return lo.Map(c.AllRel(), func(rel parser.IRelContext, _ int) *model.Relation {
		return rel.Accept(v).(*model.Relation)
	})
}

func (v *RelationVisitor) VisitSingleRel(c *parser.SingleRelContext) interface{} {
	return &model.Relation{Direct: model.ObjectName(c.Single().ID().GetText())}
}

func (v *RelationVisitor) VisitWildcardRel(c *parser.WildcardRelContext) interface{} {
	return &model.Relation{Wildcard: model.ObjectName(c.Wildcard().ID().GetText())}
}

func (v *RelationVisitor) VisitSubjectRel(c *parser.SubjectRelContext) interface{} {
	return &model.Relation{Subject: &model.SubjectRelation{
		Object:   model.ObjectName(c.Subject().ID(0).GetText()),
		Relation: model.RelationName(c.Subject().ID(1).GetText()),
	}}
}

func TestRelationParser(t *testing.T) {
	// input := antlr.NewInputStream("user | user:* | group#member | parent->viewer")
	input := antlr.NewInputStream("user | group | user:* | group#member")
	// input, _ := antlr.NewFileStream("./parser_test.txt")
	lexer := parser.NewAzmLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	// listener := NewAzmListener()

	rTree := p.Relation()
	// antlr.ParseTreeWalkerDefault.Walk(listener, rTree)
	// rel, err := listener.GetRelation()
	// if err != nil {
	//     panic(err)
	// }

	// fmt.Printf("%v\n", rel)

	var v RelationVisitor
	rel := v.Visit(rTree).([]*model.Relation)
	assert.Len(t, rel, 4)
	assert.Equal(t, model.ObjectName("user"), rel[0].Direct)
	assert.Equal(t, model.ObjectName("group"), rel[1].Direct)
	assert.Equal(t, model.ObjectName("user"), rel[2].Wildcard)
	assert.Equal(t, model.ObjectName("group"), rel[3].Subject.Object)
	assert.Equal(t, model.RelationName("member"), rel[3].Subject.Relation)

	// pTree := p.Permission()
	// antlr.ParseTreeWalkerDefault.Walk(listener, pTree)
}

func TestPermissionParser(t *testing.T) {

}
