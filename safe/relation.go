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

// Relation identifier.
type Relation struct {
	*dsc3.Relation
	HasSubjectRelation bool
}

// Relation selector.
type Relations struct {
	*Relation
}

func NewRelation(i *dsc3.Relation) *Relation { return &Relation{i, true} }

func GetRelation(i *dsr3.GetRelationRequest) *Relation {
	return &Relation{
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

func GetRelations(i *dsr3.GetRelationsRequest) *Relations {
	subjRel := ""
	if i.SubjectRelation != nil {
		subjRel = *i.SubjectRelation
	}

	return &Relations{
		&Relation{
			Relation: &dsc3.Relation{
				ObjectType:      i.ObjectType,
				ObjectId:        i.ObjectId,
				Relation:        i.Relation,
				SubjectType:     i.SubjectType,
				SubjectId:       i.SubjectId,
				SubjectRelation: subjRel,
			},
			HasSubjectRelation: i.SubjectRelation != nil,
		},
	}
}

func (i *Relation) Object() *dsc3.ObjectIdentifier {
	return &dsc3.ObjectIdentifier{
		ObjectType: i.GetObjectType(),
		ObjectId:   i.GetObjectId(),
	}
}

func (i *Relation) Subject() *dsc3.ObjectIdentifier {
	return &dsc3.ObjectIdentifier{
		ObjectType: i.GetSubjectType(),
		ObjectId:   i.GetSubjectId(),
	}
}

func (i *Relation) Validate(mc *cache.Cache) error {
	if i == nil || i.Relation == nil {
		return derr.ErrInvalidRelation.Msg("relation not set (nil)")
	}

	if IsNotSet(i.GetRelation()) {
		return derr.ErrInvalidRelation.Msg("relation")
	}

	if err := NewObjectIdentifier(i.Object()).Validate(mc); err != nil {
		return err
	}

	if err := NewObjectIdentifier(i.Subject()).Validate(mc); err != nil {
		return err
	}

	if mc == nil {
		return nil
	}

	return mc.ValidateRelation(i.Relation)
}

func (i *Relations) Validate(mc *cache.Cache) error {
	if i == nil || i.Relation.Relation == nil {
		return derr.ErrInvalidRelation.Msg("relation not set (nil)")
	}

	if err := NewObjectIdentifier(i.Object()).Validate(mc); err != nil {
		return err
	}

	if err := NewObjectIdentifier(i.Subject()).Validate(mc); err != nil {
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

func (i *Relation) Hash() string {
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
