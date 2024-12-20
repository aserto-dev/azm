package safe

import (
	"iter"

	dsc3 "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr3 "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
)

type SafeChecks struct {
	*dsr3.ChecksRequest
}

func Checks(i *dsr3.ChecksRequest) *SafeChecks {
	return &SafeChecks{i}
}

// Checks returns an iterator that materializes all checks in order.
func (c *SafeChecks) Checks() iter.Seq[SafeCheck] {
	return func(yield func(SafeCheck) bool) {
		defaults := &dsc3.RelationIdentifier{
			ObjectType:  c.Default.ObjectType,
			ObjectId:    c.Default.ObjectId,
			Relation:    c.Default.Relation,
			SubjectType: c.Default.SubjectType,
			SubjectId:   c.Default.SubjectId,
		}

		for _, check := range c.ChecksRequest.Checks {
			req := &dsr3.CheckRequest{
				ObjectType:  fallback(check.ObjectType, defaults.ObjectType),
				ObjectId:    fallback(check.ObjectId, defaults.ObjectId),
				Relation:    fallback(check.Relation, defaults.Relation),
				SubjectType: fallback(check.SubjectType, defaults.SubjectType),
				SubjectId:   fallback(check.SubjectId, defaults.SubjectId),
			}
			if !yield(SafeCheck{req}) {
				break
			}
		}
	}
}

func fallback[T comparable](val, fallback T) T {
	var def T
	if val == def {
		return fallback
	}
	return val

}
