// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import "github.com/antlr4-go/antlr/v4"

// AzmListener is a complete listener for a parse tree produced by AzmParser.
type AzmListener interface {
	antlr.ParseTreeListener

	// EnterProg is called when entering the prog production.
	EnterProg(c *ProgContext)

	// EnterStat is called when entering the stat production.
	EnterStat(c *StatContext)

	// EnterUnionRel is called when entering the unionRel production.
	EnterUnionRel(c *UnionRelContext)

	// EnterIntersectRel is called when entering the intersectRel production.
	EnterIntersectRel(c *IntersectRelContext)

	// EnterExclusionRel is called when entering the exclusionRel production.
	EnterExclusionRel(c *ExclusionRelContext)

	// EnterRel is called when entering the rel production.
	EnterRel(c *RelContext)

	// EnterSingleRel is called when entering the singleRel production.
	EnterSingleRel(c *SingleRelContext)

	// EnterSubjectRel is called when entering the subjectRel production.
	EnterSubjectRel(c *SubjectRelContext)

	// EnterWildcardRel is called when entering the wildcardRel production.
	EnterWildcardRel(c *WildcardRelContext)

	// EnterArrowRel is called when entering the arrowRel production.
	EnterArrowRel(c *ArrowRelContext)

	// ExitProg is called when exiting the prog production.
	ExitProg(c *ProgContext)

	// ExitStat is called when exiting the stat production.
	ExitStat(c *StatContext)

	// ExitUnionRel is called when exiting the unionRel production.
	ExitUnionRel(c *UnionRelContext)

	// ExitIntersectRel is called when exiting the intersectRel production.
	ExitIntersectRel(c *IntersectRelContext)

	// ExitExclusionRel is called when exiting the exclusionRel production.
	ExitExclusionRel(c *ExclusionRelContext)

	// ExitRel is called when exiting the rel production.
	ExitRel(c *RelContext)

	// ExitSingleRel is called when exiting the singleRel production.
	ExitSingleRel(c *SingleRelContext)

	// ExitSubjectRel is called when exiting the subjectRel production.
	ExitSubjectRel(c *SubjectRelContext)

	// ExitWildcardRel is called when exiting the wildcardRel production.
	ExitWildcardRel(c *WildcardRelContext)

	// ExitArrowRel is called when exiting the arrowRel production.
	ExitArrowRel(c *ArrowRelContext)
}
