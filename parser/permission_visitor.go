package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/types"
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
		return &types.Permission{}
	default:
		panic("PermissionVisitor can only visit permissions")
	}
}

func (v *PermissionVisitor) VisitUnionPerm(c *UnionPermContext) interface{} {
	return &types.Permission{
		Union: lo.Map(c.Union().AllPerm(), func(perm IPermContext, _ int) *types.PermissionTerm {
			return perm.Accept(v).(*types.PermissionTerm)
		}),
	}
}

func (v *PermissionVisitor) VisitIntersectionPerm(c *IntersectionPermContext) interface{} {
	return &types.Permission{
		Intersection: lo.Map(c.Intersection().AllPerm(), func(perm IPermContext, _ int) *types.PermissionTerm {
			return perm.Accept(v).(*types.PermissionTerm)
		}),
	}
}

func (v *PermissionVisitor) VisitExclusionPerm(c *ExclusionPermContext) interface{} {
	return &types.Permission{
		Exclusion: &types.ExclusionPermission{
			Include: c.Exclusion().Perm(0).Accept(v).(*types.PermissionTerm),
			Exclude: c.Exclusion().Perm(1).Accept(v).(*types.PermissionTerm),
		},
	}
}

func (v *PermissionVisitor) VisitDirectPerm(c *DirectPermContext) interface{} {
	return &types.PermissionTerm{RelOrPerm: types.RelationName(c.Direct().ID().GetText())}
}

func (v *PermissionVisitor) VisitArrowPerm(c *ArrowPermContext) interface{} {
	return &types.PermissionTerm{Base: types.RelationName(c.Arrow().ID(0).GetText()), RelOrPerm: types.RelationName(c.Arrow().ID(1).GetText())}
}
