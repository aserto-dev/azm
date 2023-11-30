// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import "github.com/antlr4-go/antlr/v4"

// BaseAzmListener is a complete listener for a parse tree produced by AzmParser.
type BaseAzmListener struct{}

var _ AzmListener = &BaseAzmListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseAzmListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseAzmListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseAzmListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseAzmListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterRelation is called when production relation is entered.
func (s *BaseAzmListener) EnterRelation(ctx *RelationContext) {}

// ExitRelation is called when production relation is exited.
func (s *BaseAzmListener) ExitRelation(ctx *RelationContext) {}

// EnterToUnionPerm is called when production ToUnionPerm is entered.
func (s *BaseAzmListener) EnterToUnionPerm(ctx *ToUnionPermContext) {}

// ExitToUnionPerm is called when production ToUnionPerm is exited.
func (s *BaseAzmListener) ExitToUnionPerm(ctx *ToUnionPermContext) {}

// EnterToIntersectionPerm is called when production ToIntersectionPerm is entered.
func (s *BaseAzmListener) EnterToIntersectionPerm(ctx *ToIntersectionPermContext) {}

// ExitToIntersectionPerm is called when production ToIntersectionPerm is exited.
func (s *BaseAzmListener) ExitToIntersectionPerm(ctx *ToIntersectionPermContext) {}

// EnterToExclusionPerm is called when production ToExclusionPerm is entered.
func (s *BaseAzmListener) EnterToExclusionPerm(ctx *ToExclusionPermContext) {}

// ExitToExclusionPerm is called when production ToExclusionPerm is exited.
func (s *BaseAzmListener) ExitToExclusionPerm(ctx *ToExclusionPermContext) {}

// EnterUnionRel is called when production unionRel is entered.
func (s *BaseAzmListener) EnterUnionRel(ctx *UnionRelContext) {}

// ExitUnionRel is called when production unionRel is exited.
func (s *BaseAzmListener) ExitUnionRel(ctx *UnionRelContext) {}

// EnterUnionPerm is called when production unionPerm is entered.
func (s *BaseAzmListener) EnterUnionPerm(ctx *UnionPermContext) {}

// ExitUnionPerm is called when production unionPerm is exited.
func (s *BaseAzmListener) ExitUnionPerm(ctx *UnionPermContext) {}

// EnterIntersectionPerm is called when production intersectionPerm is entered.
func (s *BaseAzmListener) EnterIntersectionPerm(ctx *IntersectionPermContext) {}

// ExitIntersectionPerm is called when production intersectionPerm is exited.
func (s *BaseAzmListener) ExitIntersectionPerm(ctx *IntersectionPermContext) {}

// EnterExclusionPerm is called when production exclusionPerm is entered.
func (s *BaseAzmListener) EnterExclusionPerm(ctx *ExclusionPermContext) {}

// ExitExclusionPerm is called when production exclusionPerm is exited.
func (s *BaseAzmListener) ExitExclusionPerm(ctx *ExclusionPermContext) {}

// EnterToSingleRel is called when production ToSingleRel is entered.
func (s *BaseAzmListener) EnterToSingleRel(ctx *ToSingleRelContext) {}

// ExitToSingleRel is called when production ToSingleRel is exited.
func (s *BaseAzmListener) ExitToSingleRel(ctx *ToSingleRelContext) {}

// EnterToWildcardRel is called when production ToWildcardRel is entered.
func (s *BaseAzmListener) EnterToWildcardRel(ctx *ToWildcardRelContext) {}

// ExitToWildcardRel is called when production ToWildcardRel is exited.
func (s *BaseAzmListener) ExitToWildcardRel(ctx *ToWildcardRelContext) {}

// EnterToSubjectRel is called when production ToSubjectRel is entered.
func (s *BaseAzmListener) EnterToSubjectRel(ctx *ToSubjectRelContext) {}

// ExitToSubjectRel is called when production ToSubjectRel is exited.
func (s *BaseAzmListener) ExitToSubjectRel(ctx *ToSubjectRelContext) {}

// EnterToArrowRel is called when production ToArrowRel is entered.
func (s *BaseAzmListener) EnterToArrowRel(ctx *ToArrowRelContext) {}

// ExitToArrowRel is called when production ToArrowRel is exited.
func (s *BaseAzmListener) ExitToArrowRel(ctx *ToArrowRelContext) {}

// EnterToSinglePerm is called when production ToSinglePerm is entered.
func (s *BaseAzmListener) EnterToSinglePerm(ctx *ToSinglePermContext) {}

// ExitToSinglePerm is called when production ToSinglePerm is exited.
func (s *BaseAzmListener) ExitToSinglePerm(ctx *ToSinglePermContext) {}

// EnterToArrowPerm is called when production ToArrowPerm is entered.
func (s *BaseAzmListener) EnterToArrowPerm(ctx *ToArrowPermContext) {}

// ExitToArrowPerm is called when production ToArrowPerm is exited.
func (s *BaseAzmListener) ExitToArrowPerm(ctx *ToArrowPermContext) {}

// EnterSingleRel is called when production singleRel is entered.
func (s *BaseAzmListener) EnterSingleRel(ctx *SingleRelContext) {}

// ExitSingleRel is called when production singleRel is exited.
func (s *BaseAzmListener) ExitSingleRel(ctx *SingleRelContext) {}

// EnterSubjectRel is called when production subjectRel is entered.
func (s *BaseAzmListener) EnterSubjectRel(ctx *SubjectRelContext) {}

// ExitSubjectRel is called when production subjectRel is exited.
func (s *BaseAzmListener) ExitSubjectRel(ctx *SubjectRelContext) {}

// EnterWildcardRel is called when production wildcardRel is entered.
func (s *BaseAzmListener) EnterWildcardRel(ctx *WildcardRelContext) {}

// ExitWildcardRel is called when production wildcardRel is exited.
func (s *BaseAzmListener) ExitWildcardRel(ctx *WildcardRelContext) {}

// EnterArrowRel is called when production arrowRel is entered.
func (s *BaseAzmListener) EnterArrowRel(ctx *ArrowRelContext) {}

// ExitArrowRel is called when production arrowRel is exited.
func (s *BaseAzmListener) ExitArrowRel(ctx *ArrowRelContext) {}
