package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type PermissionVisitor struct {
	BaseAzmVisitor
}

func (v *PermissionVisitor) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case *UnionPermContext, *IntersectionPermContext, *ExclusionPermContext:
		return t.Accept(v)
	default:
		panic("PermissionVisitor can only visit permissions")
	}
}

func (v *PermissionVisitor) VisitUnionPerm(c *UnionPermContext) interface{} {
	return &model.Permission{
		Union: lo.Map(c.Union().AllPerm(), func(perm IPermContext, _ int) *model.RelationRef {
			return perm.Accept(v).(*model.RelationRef)
		}),
	}
}

func (v *PermissionVisitor) VisitIntersectionPerm(c *IntersectionPermContext) interface{} {
	return &model.Permission{
		Intersection: lo.Map(c.Intersection().AllPerm(), func(perm IPermContext, _ int) *model.RelationRef {
			return perm.Accept(v).(*model.RelationRef)
		}),
	}
}

func (v *PermissionVisitor) VisitExclusionPerm(c *ExclusionPermContext) interface{} {
	return &model.Permission{
		Exclusion: &model.ExclusionPermission{
			Base:     c.Exclusion().Perm(0).Accept(v).(*model.RelationRef),
			Subtract: c.Exclusion().Perm(1).Accept(v).(*model.RelationRef),
		},
	}
}

func (v *PermissionVisitor) VisitSinglePerm(c *SinglePermContext) interface{} {
	return &model.RelationRef{Relation: c.Single().ID().GetText()}
}

func (v *PermissionVisitor) VisitArrowPerm(c *ArrowPermContext) interface{} {
	return &model.RelationRef{Base: model.RelationName(c.Arrow().ID(0).GetText()), Relation: c.Arrow().ID(1).GetText()}
}
