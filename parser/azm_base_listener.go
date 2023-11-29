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

// EnterProg is called when production prog is entered.
func (s *BaseAzmListener) EnterProg(ctx *ProgContext) {}

// ExitProg is called when production prog is exited.
func (s *BaseAzmListener) ExitProg(ctx *ProgContext) {}

// EnterStat is called when production stat is entered.
func (s *BaseAzmListener) EnterStat(ctx *StatContext) {}

// ExitStat is called when production stat is exited.
func (s *BaseAzmListener) ExitStat(ctx *StatContext) {}

// EnterUnionRel is called when production unionRel is entered.
func (s *BaseAzmListener) EnterUnionRel(ctx *UnionRelContext) {}

// ExitUnionRel is called when production unionRel is exited.
func (s *BaseAzmListener) ExitUnionRel(ctx *UnionRelContext) {}

// EnterIntersectRel is called when production intersectRel is entered.
func (s *BaseAzmListener) EnterIntersectRel(ctx *IntersectRelContext) {}

// ExitIntersectRel is called when production intersectRel is exited.
func (s *BaseAzmListener) ExitIntersectRel(ctx *IntersectRelContext) {}

// EnterExclusionRel is called when production exclusionRel is entered.
func (s *BaseAzmListener) EnterExclusionRel(ctx *ExclusionRelContext) {}

// ExitExclusionRel is called when production exclusionRel is exited.
func (s *BaseAzmListener) ExitExclusionRel(ctx *ExclusionRelContext) {}

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
