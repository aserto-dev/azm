package safe

import (
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/aserto-dev/azm/cache"
	"github.com/aserto-dev/azm/model"
	cerr "github.com/aserto-dev/errors"
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

type SafeObject struct {
	*dsc3.Object
}

func Object(i *dsc3.Object) *SafeObject { return &SafeObject{i} }

func (i *SafeObject) Validate(mc *cache.Cache) error {
	if i.Properties == nil {
		i.Properties = pb.NewStruct()
	}

	if mc != nil && !mc.ObjectExists(model.ObjectName(i.Type)) {
		return derr.ErrObjectTypeNotFound.Msg(i.Type)
	}

	return nil
}

func (i *SafeObject) Hash() string {
	h := fnv.New64a()
	h.Reset()

	if i.Properties != nil {
		v := i.Properties.AsMap()

		if _, err := hashstructure.Hash(
			v,
			hashstructure.FormatV2,
			&hashstructure.HashOptions{
				Hasher: h,
			},
		); err != nil {
			return DefaultHash
		}
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

type SafeObjectIdentifier struct {
	*dsc3.ObjectIdentifier
}

func ObjectIdentifier(i *dsc3.ObjectIdentifier) *SafeObjectIdentifier {
	return &SafeObjectIdentifier{i}
}

func (i *SafeObjectIdentifier) Validate(mc *cache.Cache) *cerr.AsertoError {
	if i.ObjectIdentifier == nil {
		return derr.ErrInvalidObjectIdentifier.Msg(objectIdentifierNil)
	}

	// #1 check is type field is set.
	if IsNotSet(i.GetType()) {
		return derr.ErrInvalidObjectIdentifier.Msg(objectIdentifierType)
	}

	// #2 check if id field is set.
	if IsNotSet(i.GetId()) {
		return derr.ErrInvalidObjectIdentifier.Msg(objectIdentifierID)
	}

	// #3 check if type exists.
	if mc != nil && !mc.ObjectExists(model.ObjectName(i.Type)) {
		return derr.ErrObjectTypeNotFound.Msg(i.Type)
	}

	return nil
}

func (i *SafeObjectIdentifier) Equal(n *dsc3.ObjectIdentifier) bool {
	return strings.EqualFold(i.GetId(), n.GetId()) && strings.EqualFold(i.GetType(), n.GetType())
}

func (i *SafeObjectIdentifier) IsComplete() bool {
	return i != nil && i.GetType() != "" && i.GetId() != ""
}

type SafeObjectSelector struct {
	*dsc3.ObjectIdentifier
}

func ObjectSelector(i *dsc3.ObjectIdentifier) *SafeObjectSelector { return &SafeObjectSelector{i} }

// Validate rules:
// valid states
// - empty object
// - type only
// - type + key.
func (i *SafeObjectSelector) Validate(mc *cache.Cache) error {
	// nil not allowed
	if i.ObjectIdentifier == nil {
		return derr.ErrInvalidObjectSelector.Msg(objectIdentifierNil)
	}

	switch {
	case IsSet(i.GetType()):
		// check if type exists.
		if mc != nil && !mc.ObjectExists(model.ObjectName(i.Type)) {
			return derr.ErrObjectTypeNotFound.Msg(i.Type)
		}
	case IsSet(i.GetId()):
		// can't have id without type.
		return derr.ErrInvalidObjectSelector.Msg(objectIdentifierType)
	}

	return nil
}

func (i *SafeObjectSelector) IsComplete() bool {
	return IsSet(i.GetType()) && IsSet(i.GetId())
}
