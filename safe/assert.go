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
			ObjectType:  m.CheckRelation.ObjectType,
			ObjectId:    m.CheckRelation.ObjectId,
			Relation:    m.CheckRelation.Relation,
			SubjectType: m.CheckRelation.SubjectType,
			SubjectId:   m.CheckRelation.SubjectId,
		}
		return Check(c).Validate(mc)
	case *dsa3.Assert_CheckPermission:
		c := &dsr3.CheckRequest{
			ObjectType:  m.CheckPermission.ObjectType,
			ObjectId:    m.CheckPermission.ObjectId,
			Relation:    m.CheckPermission.Permission,
			SubjectType: m.CheckPermission.SubjectType,
			SubjectId:   m.CheckPermission.SubjectId,
		}
		return Check(c).Validate(mc)
	default:
		return derr.ErrInvalidRequest.Msg("msg")
	}
}
