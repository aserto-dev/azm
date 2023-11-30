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

	// EnterUnionPerm is called when entering the unionPerm production.
	EnterUnionPerm(c *UnionPermContext)

	// EnterIntersectionPerm is called when entering the intersectionPerm production.
	EnterIntersectionPerm(c *IntersectionPermContext)

	// EnterExclusionPerm is called when entering the exclusionPerm production.
	EnterExclusionPerm(c *ExclusionPermContext)

	// EnterSingleRel is called when entering the SingleRel production.
	EnterSingleRel(c *SingleRelContext)

	// EnterWildcardRel is called when entering the WildcardRel production.
	EnterWildcardRel(c *WildcardRelContext)

	// EnterSubjectRel is called when entering the SubjectRel production.
	EnterSubjectRel(c *SubjectRelContext)

	// EnterArrowRel is called when entering the ArrowRel production.
	EnterArrowRel(c *ArrowRelContext)

	// EnterSinglePerm is called when entering the SinglePerm production.
	EnterSinglePerm(c *SinglePermContext)

	// EnterArrowPerm is called when entering the ArrowPerm production.
	EnterArrowPerm(c *ArrowPermContext)

	// EnterSingle is called when entering the single production.
	EnterSingle(c *SingleContext)

	// EnterSubject is called when entering the subject production.
	EnterSubject(c *SubjectContext)

	// EnterWildcard is called when entering the wildcard production.
	EnterWildcard(c *WildcardContext)

	// EnterArrow is called when entering the arrow production.
	EnterArrow(c *ArrowContext)

	// ExitRelation is called when exiting the relation production.
	ExitRelation(c *RelationContext)

	// ExitPermission is called when exiting the permission production.
	ExitPermission(c *PermissionContext)

	// ExitUnionPerm is called when exiting the unionPerm production.
	ExitUnionPerm(c *UnionPermContext)

	// ExitIntersectionPerm is called when exiting the intersectionPerm production.
	ExitIntersectionPerm(c *IntersectionPermContext)

	// ExitExclusionPerm is called when exiting the exclusionPerm production.
	ExitExclusionPerm(c *ExclusionPermContext)

	// ExitSingleRel is called when exiting the SingleRel production.
	ExitSingleRel(c *SingleRelContext)

	// ExitWildcardRel is called when exiting the WildcardRel production.
	ExitWildcardRel(c *WildcardRelContext)

	// ExitSubjectRel is called when exiting the SubjectRel production.
	ExitSubjectRel(c *SubjectRelContext)

	// ExitArrowRel is called when exiting the ArrowRel production.
	ExitArrowRel(c *ArrowRelContext)

	// ExitSinglePerm is called when exiting the SinglePerm production.
	ExitSinglePerm(c *SinglePermContext)

	// ExitArrowPerm is called when exiting the ArrowPerm production.
	ExitArrowPerm(c *ArrowPermContext)

	// ExitSingle is called when exiting the single production.
	ExitSingle(c *SingleContext)

	// ExitSubject is called when exiting the subject production.
	ExitSubject(c *SubjectContext)

	// ExitWildcard is called when exiting the wildcard production.
	ExitWildcard(c *WildcardContext)

	// ExitArrow is called when exiting the arrow production.
	ExitArrow(c *ArrowContext)
}
