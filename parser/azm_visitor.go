// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by AzmParser.
type AzmVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by AzmParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by AzmParser#ToUnionPerm.
	VisitToUnionPerm(ctx *ToUnionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#ToIntersectionPerm.
	VisitToIntersectionPerm(ctx *ToIntersectionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#ToExclusionPerm.
	VisitToExclusionPerm(ctx *ToExclusionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#unionRel.
	VisitUnionRel(ctx *UnionRelContext) interface{}

	// Visit a parse tree produced by AzmParser#unionPerm.
	VisitUnionPerm(ctx *UnionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#intersectionPerm.
	VisitIntersectionPerm(ctx *IntersectionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#exclusionPerm.
	VisitExclusionPerm(ctx *ExclusionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#ToSingleRel.
	VisitToSingleRel(ctx *ToSingleRelContext) interface{}

	// Visit a parse tree produced by AzmParser#ToWildcardRel.
	VisitToWildcardRel(ctx *ToWildcardRelContext) interface{}

	// Visit a parse tree produced by AzmParser#ToSubjectRel.
	VisitToSubjectRel(ctx *ToSubjectRelContext) interface{}

	// Visit a parse tree produced by AzmParser#ToArrowRel.
	VisitToArrowRel(ctx *ToArrowRelContext) interface{}

	// Visit a parse tree produced by AzmParser#ToSinglePerm.
	VisitToSinglePerm(ctx *ToSinglePermContext) interface{}

	// Visit a parse tree produced by AzmParser#ToArrowPerm.
	VisitToArrowPerm(ctx *ToArrowPermContext) interface{}

	// Visit a parse tree produced by AzmParser#singleRel.
	VisitSingleRel(ctx *SingleRelContext) interface{}

	// Visit a parse tree produced by AzmParser#subjectRel.
	VisitSubjectRel(ctx *SubjectRelContext) interface{}

	// Visit a parse tree produced by AzmParser#wildcardRel.
	VisitWildcardRel(ctx *WildcardRelContext) interface{}

	// Visit a parse tree produced by AzmParser#arrowRel.
	VisitArrowRel(ctx *ArrowRelContext) interface{}
}
