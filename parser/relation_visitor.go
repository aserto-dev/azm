package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/samber/lo"
)

type RelationVisitor struct {
	BaseAzmVisitor
}

func (v *RelationVisitor) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case *RelationContext:
		return t.Accept(v)
	default:
		panic("RelationVisitor can only visit relations")
	}
}

func (v *RelationVisitor) VisitRelation(c *RelationContext) interface{} {
	return lo.Map(c.AllRel(), func(rel IRelContext, _ int) *RelationRef {
		if term, ok := rel.Accept(v).(*RelationRef); ok {
			return term
		}

		return &RelationRef{}
	})
}

func (v *RelationVisitor) VisitDirectRel(c *DirectRelContext) interface{} {
	return &RelationRef{Object: ObjectName(c.Direct().ID().GetText())}
}

func (v *RelationVisitor) VisitWildcardRel(c *WildcardRelContext) interface{} {
	return &RelationRef{
		Object:   ObjectName(c.Wildcard().ID().GetText()),
		Relation: "*",
	}
}

func (v *RelationVisitor) VisitSubjectRel(c *SubjectRelContext) interface{} {
	return &RelationRef{
		Object:   ObjectName(c.Subject().ID(0).GetText()),
		Relation: RelationName(c.Subject().ID(1).GetText()),
	}
}
