package safe

import (
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/aserto-dev/azm/cache"
	"github.com/aserto-dev/azm/model"
	dsc3 "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/aserto-dev/go-directory/pkg/pb"

	"github.com/mitchellh/hashstructure/v2"
)

const DefaultHash string = `0`

const (
	objectIdentifierNil  string = "not set (nil)"
	objectIdentifierType string = "type"
	objectIdentifierID   string = "id"
)

func IsSet(s string) bool {
	return strings.TrimSpace(s) != ""
}

func IsNotSet(s string) bool {
	return strings.TrimSpace(s) == ""
}

type Object struct {
	*dsc3.Object
}

func NewObject(i *dsc3.Object) *Object { return &Object{i} }

func (i *Object) Validate(mc *cache.Cache) error {
	if i.Properties == nil {
		i.Properties = pb.NewStruct()
	}

	if mc != nil && !mc.ObjectExists(model.ObjectName(i.Object.Type)) {
		return derr.ErrObjectTypeNotFound.Msg(i.Object.Type)
	}

	return nil
}

func (i *Object) Hash() string {
	h := fnv.New64a()
	h.Reset()

	if i.Properties != nil {
		v := i.Properties.AsMap()
		hash, err := hashstructure.Hash(
			v,
			hashstructure.FormatV2,
			&hashstructure.HashOptions{
				Hasher: h,
			},
		)
		if err != nil {
			return DefaultHash
		}
		_ = hash
	}

	if _, err := h.Write([]byte(i.GetType())); err != nil {
		return DefaultHash
	}
	if _, err := h.Write([]byte(i.GetId())); err != nil {
		return DefaultHash
	}

	if _, err := h.Write([]byte(i.GetDisplayName())); err != nil {
		return DefaultHash
	}

	return strconv.FormatUint(h.Sum64(), 10)
}

type ObjectIdentifier struct {
	*dsc3.ObjectIdentifier
}

func NewObjectIdentifier(i *dsc3.ObjectIdentifier) *ObjectIdentifier { return &ObjectIdentifier{i} }

func (i *ObjectIdentifier) Validate(mc *cache.Cache) error {
	if i.ObjectIdentifier == nil {
		return derr.ErrInvalidObjectIdentifier.Msg(objectIdentifierNil)
	}

	// #1 check is type field is set.
	if IsNotSet(i.GetObjectType()) {
		return derr.ErrInvalidObjectIdentifier.Msg(objectIdentifierType)
	}

	// #2 check if id field is set.
	if IsNotSet(i.GetObjectId()) {
		return derr.ErrInvalidObjectIdentifier.Msg(objectIdentifierID)
	}

	// #3 check if type exists.
	if mc != nil && !mc.ObjectExists(model.ObjectName(i.ObjectIdentifier.ObjectType)) {
		return derr.ErrObjectTypeNotFound.Msg(i.ObjectIdentifier.ObjectType)
	}

	return nil
}

func (i *ObjectIdentifier) Equal(n *dsc3.ObjectIdentifier) bool {
	return strings.EqualFold(i.ObjectIdentifier.GetObjectId(), n.GetObjectId()) && strings.EqualFold(i.ObjectIdentifier.GetObjectType(), n.GetObjectType())
}

func (i *ObjectIdentifier) IsComplete() bool {
	return i != nil && i.GetObjectType() != "" && i.GetObjectId() != ""
}

type ObjectSelector struct {
	*dsc3.ObjectIdentifier
}

func NewObjectSelector(i *dsc3.ObjectIdentifier) *ObjectSelector { return &ObjectSelector{i} }

// Validate rules:
// valid states
// - empty object
// - type only
// - type + key.
func (i *ObjectSelector) Validate(mc *cache.Cache) error {
	// nil not allowed
	if i.ObjectIdentifier == nil {
		return derr.ErrInvalidObjectSelector.Msg(objectIdentifierNil)
	}

	switch {
	case IsSet(i.GetObjectType()):
		// check if type exists.
		if mc != nil && !mc.ObjectExists(model.ObjectName(i.ObjectIdentifier.ObjectType)) {
			return derr.ErrObjectTypeNotFound.Msg(i.ObjectIdentifier.ObjectType)
		}
	case IsSet(i.GetObjectId()):
		// can't have id without type.
		return derr.ErrInvalidObjectSelector.Msg(objectIdentifierType)
	}

	return nil
}

func (i *ObjectSelector) IsComplete() bool {
	return IsSet(i.GetObjectType()) && IsSet(i.GetObjectId())
}
