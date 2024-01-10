package safe

import (
	"hash/fnv"
	"strconv"

	"github.com/aserto-dev/azm/cache"
	"github.com/aserto-dev/azm/model"
	dsc3 "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr3 "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
)

// SafeRelation identifier.
type SafeRelation struct {
	*dsc3.Relation
	HasSubjectRelation bool
}

func Relation(i *dsc3.Relation) *SafeRelation { return &SafeRelation{i, true} }

type SafeRelationIdentifier struct {
	*model.RelationRef
}

func RelationIdentifier(i *model.RelationRef) *SafeRelationIdentifier {
	return &SafeRelationIdentifier{i}
}

// Relation selector.
type SafeRelations struct {
	*SafeRelation
}

func GetRelation(i *dsr3.GetRelationRequest) *SafeRelation {
	return &SafeRelation{
		Relation: &dsc3.Relation{
			ObjectType:      i.ObjectType,
			ObjectId:        i.ObjectId,
			Relation:        i.Relation,
			SubjectType:     i.SubjectType,
			SubjectId:       i.SubjectId,
			SubjectRelation: i.SubjectRelation,
		},
		HasSubjectRelation: true,
	}
}

func GetRelations(i *dsr3.GetRelationsRequest) *SafeRelations {
	return &SafeRelations{
		&SafeRelation{
			Relation: &dsc3.Relation{
				ObjectType:      i.ObjectType,
				ObjectId:        i.ObjectId,
				Relation:        i.Relation,
				SubjectType:     i.SubjectType,
				SubjectId:       i.SubjectId,
				SubjectRelation: i.SubjectRelation,
			},
			HasSubjectRelation: i.SubjectRelation != "" || i.WithEmptySubjectRelation,
		},
	}
}

func (i *SafeRelation) Object() *dsc3.ObjectIdentifier {
	return &dsc3.ObjectIdentifier{
		ObjectType: i.GetObjectType(),
		ObjectId:   i.GetObjectId(),
	}
}

func (i *SafeRelation) Subject() *dsc3.ObjectIdentifier {
	return &dsc3.ObjectIdentifier{
		ObjectType: i.GetSubjectType(),
		ObjectId:   i.GetSubjectId(),
	}
}

func (i *SafeRelation) Validate(mc *cache.Cache) error {
	if i == nil || i.Relation == nil {
		return derr.ErrInvalidRelation.Msg("relation not set (nil)")
	}

	if IsNotSet(i.GetRelation()) {
		return derr.ErrInvalidRelation.Msg("relation")
	}

	if err := ObjectIdentifier(i.Object()).Validate(mc); err != nil {
		return err
	}

	if err := ObjectIdentifier(i.Subject()).Validate(mc); err != nil {
		return err
	}

	if mc == nil {
		return nil
	}

	return mc.ValidateRelation(i.Relation)
}

func (i *SafeRelations) Validate(mc *cache.Cache) error {
	if i == nil || i.SafeRelation.Relation == nil {
		return derr.ErrInvalidRelation.Msg("relation not set (nil)")
	}

	if err := ObjectSelector(i.Object()).Validate(mc); err != nil {
		return err
	}

	if err := ObjectSelector(i.Subject()).Validate(mc); err != nil {
		return err
	}

	if IsSet(i.GetRelation()) {
		if IsNotSet(i.GetObjectType()) {
			return derr.ErrInvalidRelation.Msg("object type not set")
		}

		if mc != nil && !mc.RelationExists(model.ObjectName(i.GetObjectType()), model.RelationName(i.GetRelation())) {
			return derr.ErrRelationNotFound.Msg(i.GetObjectType() + ":" + i.GetRelation())
		}
	}

	if IsSet(i.GetSubjectRelation()) {
		if IsNotSet(i.GetSubjectType()) {
			return derr.ErrInvalidRelation.Msg("subject type not set")
		}

		if mc != nil && !mc.RelationExists(model.ObjectName(i.GetSubjectType()), model.RelationName(i.GetSubjectRelation())) {
			return derr.ErrRelationNotFound.Msg(i.GetSubjectType() + ":" + i.GetSubjectRelation())
		}
	}

	return nil
}

func (i *SafeRelation) Hash() string {
	h := fnv.New64a()
	h.Reset()

	if i != nil && i.Relation != nil {
		if _, err := h.Write([]byte(i.GetObjectId())); err != nil {
			return DefaultHash
		}
		if _, err := h.Write([]byte(i.GetObjectType())); err != nil {
			return DefaultHash
		}
		if _, err := h.Write([]byte(i.GetRelation())); err != nil {
			return DefaultHash
		}
		if _, err := h.Write([]byte(i.GetSubjectId())); err != nil {
			return DefaultHash
		}
		if _, err := h.Write([]byte(i.GetSubjectType())); err != nil {
			return DefaultHash
		}
		if _, err := h.Write([]byte(i.GetSubjectRelation())); err != nil {
			return DefaultHash
		}
	}

	return strconv.FormatUint(h.Sum64(), 10)
}

type RelationScope int

const (
	AsRelation RelationScope = iota
	AsPermission
	AsEither
)

func (r *SafeRelationIdentifier) Validate(scope RelationScope, mc *cache.Cache) error {
	if r == nil || r.RelationRef == nil {
		return derr.ErrInvalidRelation.Msg("relation not set (nil)")
	}

	if r.Object == "" {
		return derr.ErrInvalidRelation.Msg("object")
	}

	if r.Relation == "" {
		return derr.ErrInvalidRelation.Msg("relation")
	}

	if mc == nil {
		return nil
	}

	switch scope {
	case AsRelation:
		if !mc.RelationExists(r.Object, r.Relation) {
			return derr.ErrRelationNotFound.Msgf("relation: %s", r)
		}
	case AsPermission:
		if !mc.PermissionExists(r.Object, r.Relation) {
			return derr.ErrPermissionNotFound.Msgf("permission: %s", r)
		}
	case AsEither:
		if !mc.RelationExists(r.Object, r.Relation) && !mc.PermissionExists(r.Object, r.Relation) {
			return derr.ErrRelationNotFound.Msgf("relation: %s", r)
		}
	}

	return nil
}
