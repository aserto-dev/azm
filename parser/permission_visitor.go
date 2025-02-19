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
	case *PermissionContext:
		return &model.Permission{}
	default:
		panic("PermissionVisitor can only visit permissions")
	}
}

func (v *PermissionVisitor) VisitUnionPerm(c *UnionPermContext) interface{} {
	return &model.Permission{
		Union: lo.Map(c.Union().AllPerm(), func(perm IPermContext, _ int) *model.PermissionTerm {
			if term, ok := perm.Accept(v).(*model.PermissionTerm); ok {
				return term
			}

			return &model.PermissionTerm{}
		}),
	}
}

func (v *PermissionVisitor) VisitIntersectionPerm(c *IntersectionPermContext) interface{} {
	return &model.Permission{
		Intersection: lo.Map(c.Intersection().AllPerm(), func(perm IPermContext, _ int) *model.PermissionTerm {
			return perm.Accept(v).(*model.PermissionTerm)
		}),
	}
}

func (v *PermissionVisitor) VisitExclusionPerm(c *ExclusionPermContext) interface{} {
	return &model.Permission{
		Exclusion: &model.ExclusionPermission{
			Include: c.Exclusion().Perm(0).Accept(v).(*model.PermissionTerm),
			Exclude: c.Exclusion().Perm(1).Accept(v).(*model.PermissionTerm),
		},
	}
}

func (v *PermissionVisitor) VisitDirectPerm(c *DirectPermContext) interface{} {
	return &model.PermissionTerm{RelOrPerm: model.RelationName(c.Direct().ID().GetText())}
}

func (v *PermissionVisitor) VisitArrowPerm(c *ArrowPermContext) interface{} {
	return &model.PermissionTerm{Base: model.RelationName(c.Arrow().ID(0).GetText()), RelOrPerm: model.RelationName(c.Arrow().ID(1).GetText())}
}
