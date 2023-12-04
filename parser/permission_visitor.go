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
		Union: lo.Map(c.Union().AllPerm(), func(perm IPermContext, _ int) *model.PermissionRef {
			return perm.Accept(v).(*model.PermissionRef)
		}),
	}
}

func (v *PermissionVisitor) VisitIntersectionPerm(c *IntersectionPermContext) interface{} {
	return &model.Permission{
		Intersection: lo.Map(c.Intersection().AllPerm(), func(perm IPermContext, _ int) *model.PermissionRef {
			return perm.Accept(v).(*model.PermissionRef)
		}),
	}
}

func (v *PermissionVisitor) VisitExclusionPerm(c *ExclusionPermContext) interface{} {
	return &model.Permission{
		Exclusion: &model.ExclusionPermission{
			Include: c.Exclusion().Perm(0).Accept(v).(*model.PermissionRef),
			Exclude: c.Exclusion().Perm(1).Accept(v).(*model.PermissionRef),
		},
	}
}

func (v *PermissionVisitor) VisitDirectPerm(c *DirectPermContext) interface{} {
	return &model.PermissionRef{RelOrPerm: c.Direct().ID().GetText()}
}

func (v *PermissionVisitor) VisitArrowPerm(c *ArrowPermContext) interface{} {
	return &model.PermissionRef{Base: model.RelationName(c.Arrow().ID(0).GetText()), RelOrPerm: c.Arrow().ID(1).GetText()}
}
