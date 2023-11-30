// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import "github.com/antlr4-go/antlr/v4"

// AzmListener is a complete listener for a parse tree produced by AzmParser.
type AzmListener interface {
	antlr.ParseTreeListener

	// EnterRelation is called when entering the relation production.
	EnterRelation(c *RelationContext)

	// EnterPermission is called when entering the permission production.
	EnterPermission(c *PermissionContext)

	// EnterUnion is called when entering the union production.
	EnterUnion(c *UnionContext)

	// EnterIntersection is called when entering the intersection production.
	EnterIntersection(c *IntersectionContext)

	// EnterExclusion is called when entering the exclusion production.
	EnterExclusion(c *ExclusionContext)

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

	// ExitRelation is called when exiting the relation production.
	ExitRelation(c *RelationContext)

	// ExitPermission is called when exiting the permission production.
	ExitPermission(c *PermissionContext)

	// ExitUnion is called when exiting the union production.
	ExitUnion(c *UnionContext)

	// ExitIntersection is called when exiting the intersection production.
	ExitIntersection(c *IntersectionContext)

	// ExitExclusion is called when exiting the exclusion production.
	ExitExclusion(c *ExclusionContext)

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
