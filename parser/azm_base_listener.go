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

// EnterUnion is called when production union is entered.
func (s *BaseAzmListener) EnterUnion(ctx *UnionContext) {}

// ExitUnion is called when production union is exited.
func (s *BaseAzmListener) ExitUnion(ctx *UnionContext) {}

// EnterIntersection is called when production intersection is entered.
func (s *BaseAzmListener) EnterIntersection(ctx *IntersectionContext) {}

// ExitIntersection is called when production intersection is exited.
func (s *BaseAzmListener) ExitIntersection(ctx *IntersectionContext) {}

// EnterExclusion is called when production exclusion is entered.
func (s *BaseAzmListener) EnterExclusion(ctx *ExclusionContext) {}

// ExitExclusion is called when production exclusion is exited.
func (s *BaseAzmListener) ExitExclusion(ctx *ExclusionContext) {}

// EnterRel is called when production rel is entered.
func (s *BaseAzmListener) EnterRel(ctx *RelContext) {}

// ExitRel is called when production rel is exited.
func (s *BaseAzmListener) ExitRel(ctx *RelContext) {}

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
