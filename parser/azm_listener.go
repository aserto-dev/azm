// Code generated from Azm.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Azm
import "github.com/antlr4-go/antlr/v4"

// AzmListener is a complete listener for a parse tree produced by AzmParser.
type AzmListener interface {
	antlr.ParseTreeListener

	// EnterRelation is called when entering the relation production.
	EnterRelation(c *RelationContext)

	// EnterToUnionPerm is called when entering the ToUnionPerm production.
	EnterToUnionPerm(c *ToUnionPermContext)

	// EnterToIntersectionPerm is called when entering the ToIntersectionPerm production.
	EnterToIntersectionPerm(c *ToIntersectionPermContext)

	// EnterToExclusionPerm is called when entering the ToExclusionPerm production.
	EnterToExclusionPerm(c *ToExclusionPermContext)

	// EnterUnionRel is called when entering the unionRel production.
	EnterUnionRel(c *UnionRelContext)

	// EnterUnionPerm is called when entering the unionPerm production.
	EnterUnionPerm(c *UnionPermContext)

	// EnterIntersectionPerm is called when entering the intersectionPerm production.
	EnterIntersectionPerm(c *IntersectionPermContext)

	// EnterExclusionPerm is called when entering the exclusionPerm production.
	EnterExclusionPerm(c *ExclusionPermContext)

	// EnterToSingleRel is called when entering the ToSingleRel production.
	EnterToSingleRel(c *ToSingleRelContext)

	// EnterToWildcardRel is called when entering the ToWildcardRel production.
	EnterToWildcardRel(c *ToWildcardRelContext)

	// EnterToSubjectRel is called when entering the ToSubjectRel production.
	EnterToSubjectRel(c *ToSubjectRelContext)

	// EnterToArrowRel is called when entering the ToArrowRel production.
	EnterToArrowRel(c *ToArrowRelContext)

	// EnterToSinglePerm is called when entering the ToSinglePerm production.
	EnterToSinglePerm(c *ToSinglePermContext)

	// EnterToArrowPerm is called when entering the ToArrowPerm production.
	EnterToArrowPerm(c *ToArrowPermContext)

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

	// ExitToUnionPerm is called when exiting the ToUnionPerm production.
	ExitToUnionPerm(c *ToUnionPermContext)

	// ExitToIntersectionPerm is called when exiting the ToIntersectionPerm production.
	ExitToIntersectionPerm(c *ToIntersectionPermContext)

	// ExitToExclusionPerm is called when exiting the ToExclusionPerm production.
	ExitToExclusionPerm(c *ToExclusionPermContext)

	// ExitUnionRel is called when exiting the unionRel production.
	ExitUnionRel(c *UnionRelContext)

	// ExitUnionPerm is called when exiting the unionPerm production.
	ExitUnionPerm(c *UnionPermContext)

	// ExitIntersectionPerm is called when exiting the intersectionPerm production.
	ExitIntersectionPerm(c *IntersectionPermContext)

	// ExitExclusionPerm is called when exiting the exclusionPerm production.
	ExitExclusionPerm(c *ExclusionPermContext)

	// ExitToSingleRel is called when exiting the ToSingleRel production.
	ExitToSingleRel(c *ToSingleRelContext)

	// ExitToWildcardRel is called when exiting the ToWildcardRel production.
	ExitToWildcardRel(c *ToWildcardRelContext)

	// ExitToSubjectRel is called when exiting the ToSubjectRel production.
	ExitToSubjectRel(c *ToSubjectRelContext)

	// ExitToArrowRel is called when exiting the ToArrowRel production.
	ExitToArrowRel(c *ToArrowRelContext)

	// ExitToSinglePerm is called when exiting the ToSinglePerm production.
	ExitToSinglePerm(c *ToSinglePermContext)

	// ExitToArrowPerm is called when exiting the ToArrowPerm production.
	ExitToArrowPerm(c *ToArrowPermContext)

	// ExitSingleRel is called when exiting the singleRel production.
	ExitSingleRel(c *SingleRelContext)

	// ExitSubjectRel is called when exiting the subjectRel production.
	ExitSubjectRel(c *SubjectRelContext)

	// ExitWildcardRel is called when exiting the wildcardRel production.
	ExitWildcardRel(c *WildcardRelContext)

	// ExitArrowRel is called when exiting the arrowRel production.
	ExitArrowRel(c *ArrowRelContext)
}
