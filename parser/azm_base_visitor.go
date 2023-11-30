// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import "github.com/antlr4-go/antlr/v4"

type BaseAzmVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseAzmVisitor) VisitRelation(ctx *RelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToUnionPerm(ctx *ToUnionPermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToIntersectionPerm(ctx *ToIntersectionPermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToExclusionPerm(ctx *ToExclusionPermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitUnionRel(ctx *UnionRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitUnionPerm(ctx *UnionPermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitIntersectionPerm(ctx *IntersectionPermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitExclusionPerm(ctx *ExclusionPermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToSingleRel(ctx *ToSingleRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToWildcardRel(ctx *ToWildcardRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToSubjectRel(ctx *ToSubjectRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToArrowRel(ctx *ToArrowRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToSinglePerm(ctx *ToSinglePermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitToArrowPerm(ctx *ToArrowPermContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitSingleRel(ctx *SingleRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitSubjectRel(ctx *SubjectRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitWildcardRel(ctx *WildcardRelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseAzmVisitor) VisitArrowRel(ctx *ArrowRelContext) interface{} {
	return v.VisitChildren(ctx)
}
