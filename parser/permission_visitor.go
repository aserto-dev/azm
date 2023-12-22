package parser

import (
	"github.com/antlr4-go/antlr/v4"
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
		return &Permission{}
	default:
		panic("PermissionVisitor can only visit permissions")
	}
}

func (v *PermissionVisitor) VisitUnionPerm(c *UnionPermContext) interface{} {
	return &Permission{
		Union: lo.Map(c.Union().AllPerm(), func(perm IPermContext, _ int) *PermissionTerm {
			return perm.Accept(v).(*PermissionTerm)
		}),
	}
}

func (v *PermissionVisitor) VisitIntersectionPerm(c *IntersectionPermContext) interface{} {
	return &Permission{
		Intersection: lo.Map(c.Intersection().AllPerm(), func(perm IPermContext, _ int) *PermissionTerm {
			return perm.Accept(v).(*PermissionTerm)
		}),
	}
}

func (v *PermissionVisitor) VisitExclusionPerm(c *ExclusionPermContext) interface{} {
	return &Permission{
		Exclusion: &ExclusionPermission{
			Include: c.Exclusion().Perm(0).Accept(v).(*PermissionTerm),
			Exclude: c.Exclusion().Perm(1).Accept(v).(*PermissionTerm),
		},
	}
}

func (v *PermissionVisitor) VisitDirectPerm(c *DirectPermContext) interface{} {
	return &PermissionTerm{RelOrPerm: RelationName(c.Direct().ID().GetText())}
}

func (v *PermissionVisitor) VisitArrowPerm(c *ArrowPermContext) interface{} {
	return &PermissionTerm{Base: RelationName(c.Arrow().ID(0).GetText()), RelOrPerm: RelationName(c.Arrow().ID(1).GetText())}
}
