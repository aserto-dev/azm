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

// EnterPermission is called when production permission is entered.
func (s *BaseAzmListener) EnterPermission(ctx *PermissionContext) {}

// ExitPermission is called when production permission is exited.
func (s *BaseAzmListener) ExitPermission(ctx *PermissionContext) {}

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

// EnterSingleRel is called when production SingleRel is entered.
func (s *BaseAzmListener) EnterSingleRel(ctx *SingleRelContext) {}

// ExitSingleRel is called when production SingleRel is exited.
func (s *BaseAzmListener) ExitSingleRel(ctx *SingleRelContext) {}

// EnterWildcardRel is called when production WildcardRel is entered.
func (s *BaseAzmListener) EnterWildcardRel(ctx *WildcardRelContext) {}

// ExitWildcardRel is called when production WildcardRel is exited.
func (s *BaseAzmListener) ExitWildcardRel(ctx *WildcardRelContext) {}

// EnterSubjectRel is called when production SubjectRel is entered.
func (s *BaseAzmListener) EnterSubjectRel(ctx *SubjectRelContext) {}

// ExitSubjectRel is called when production SubjectRel is exited.
func (s *BaseAzmListener) ExitSubjectRel(ctx *SubjectRelContext) {}

// EnterArrowRel is called when production ArrowRel is entered.
func (s *BaseAzmListener) EnterArrowRel(ctx *ArrowRelContext) {}

// ExitArrowRel is called when production ArrowRel is exited.
func (s *BaseAzmListener) ExitArrowRel(ctx *ArrowRelContext) {}

// EnterSinglePerm is called when production SinglePerm is entered.
func (s *BaseAzmListener) EnterSinglePerm(ctx *SinglePermContext) {}

// ExitSinglePerm is called when production SinglePerm is exited.
func (s *BaseAzmListener) ExitSinglePerm(ctx *SinglePermContext) {}

// EnterArrowPerm is called when production ArrowPerm is entered.
func (s *BaseAzmListener) EnterArrowPerm(ctx *ArrowPermContext) {}

// ExitArrowPerm is called when production ArrowPerm is exited.
func (s *BaseAzmListener) ExitArrowPerm(ctx *ArrowPermContext) {}

// EnterSingle is called when production single is entered.
func (s *BaseAzmListener) EnterSingle(ctx *SingleContext) {}

// ExitSingle is called when production single is exited.
func (s *BaseAzmListener) ExitSingle(ctx *SingleContext) {}

// EnterSubject is called when production subject is entered.
func (s *BaseAzmListener) EnterSubject(ctx *SubjectContext) {}

// ExitSubject is called when production subject is exited.
func (s *BaseAzmListener) ExitSubject(ctx *SubjectContext) {}

// EnterWildcard is called when production wildcard is entered.
func (s *BaseAzmListener) EnterWildcard(ctx *WildcardContext) {}

// ExitWildcard is called when production wildcard is exited.
func (s *BaseAzmListener) ExitWildcard(ctx *WildcardContext) {}

// EnterArrow is called when production arrow is entered.
func (s *BaseAzmListener) EnterArrow(ctx *ArrowContext) {}

// ExitArrow is called when production arrow is exited.
func (s *BaseAzmListener) ExitArrow(ctx *ArrowContext) {}
