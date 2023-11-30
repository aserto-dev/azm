// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by AzmParser.
type AzmVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by AzmParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by AzmParser#permission.
	VisitPermission(ctx *PermissionContext) interface{}

	// Visit a parse tree produced by AzmParser#unionPerm.
	VisitUnionPerm(ctx *UnionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#intersectionPerm.
	VisitIntersectionPerm(ctx *IntersectionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#exclusionPerm.
	VisitExclusionPerm(ctx *ExclusionPermContext) interface{}

	// Visit a parse tree produced by AzmParser#SingleRel.
	VisitSingleRel(ctx *SingleRelContext) interface{}

	// Visit a parse tree produced by AzmParser#WildcardRel.
	VisitWildcardRel(ctx *WildcardRelContext) interface{}

	// Visit a parse tree produced by AzmParser#SubjectRel.
	VisitSubjectRel(ctx *SubjectRelContext) interface{}

	// Visit a parse tree produced by AzmParser#ArrowRel.
	VisitArrowRel(ctx *ArrowRelContext) interface{}

	// Visit a parse tree produced by AzmParser#SinglePerm.
	VisitSinglePerm(ctx *SinglePermContext) interface{}

	// Visit a parse tree produced by AzmParser#ArrowPerm.
	VisitArrowPerm(ctx *ArrowPermContext) interface{}

	// Visit a parse tree produced by AzmParser#single.
	VisitSingle(ctx *SingleContext) interface{}

	// Visit a parse tree produced by AzmParser#subject.
	VisitSubject(ctx *SubjectContext) interface{}

	// Visit a parse tree produced by AzmParser#wildcard.
	VisitWildcard(ctx *WildcardContext) interface{}

	// Visit a parse tree produced by AzmParser#arrow.
	VisitArrow(ctx *ArrowContext) interface{}
}
