package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
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
	return lo.Map(c.AllRel(), func(rel IRelContext, _ int) *model.RelationTerm {
		return rel.Accept(v).(*model.RelationTerm)
	})
}

func (v *RelationVisitor) VisitDirectRel(c *DirectRelContext) interface{} {
	return &model.RelationTerm{Direct: model.ObjectName(c.Direct().ID().GetText())}
}

func (v *RelationVisitor) VisitWildcardRel(c *WildcardRelContext) interface{} {
	return &model.RelationTerm{Wildcard: model.ObjectName(c.Wildcard().ID().GetText())}
}

func (v *RelationVisitor) VisitSubjectRel(c *SubjectRelContext) interface{} {
	return &model.RelationTerm{Subject: &model.SubjectRelation{
		RelationRef: &model.RelationRef{
			Object:   model.ObjectName(c.Subject().ID(0).GetText()),
			Relation: model.RelationName(c.Subject().ID(1).GetText()),
		},
		SubjectTypes: []model.ObjectName{},
	}}
}
