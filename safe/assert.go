package safe

import (
	"github.com/aserto-dev/azm/cache"

	dsa3 "github.com/aserto-dev/go-directory/aserto/directory/assertion/v3"
	dsr3 "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
)

type SafeAssert struct {
	*dsa3.Assert
}

func Assert(i *dsa3.Assert) *SafeAssert {
	return &SafeAssert{i}
}

func (i *SafeAssert) Validate(mc *cache.Cache) error {
	switch m := i.Msg.(type) {
	case *dsa3.Assert_Check:
		return Check(m.Check).Validate(mc)
	case *dsa3.Assert_CheckRelation:
		c := &dsr3.CheckRequest{
			ObjectType:  m.CheckRelation.GetObjectType(),
			ObjectId:    m.CheckRelation.GetObjectId(),
			Relation:    m.CheckRelation.GetRelation(),
			SubjectType: m.CheckRelation.GetSubjectType(),
			SubjectId:   m.CheckRelation.GetSubjectId(),
		}

		return Check(c).Validate(mc)
	case *dsa3.Assert_CheckPermission:
		c := &dsr3.CheckRequest{
			ObjectType:  m.CheckPermission.GetObjectType(),
			ObjectId:    m.CheckPermission.GetObjectId(),
			Relation:    m.CheckPermission.GetPermission(),
			SubjectType: m.CheckPermission.GetSubjectType(),
			SubjectId:   m.CheckPermission.GetSubjectId(),
		}

		return Check(c).Validate(mc)
	default:
		return derr.ErrInvalidRequest.Msg("msg")
	}
}
