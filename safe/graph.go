package safe

import (
	"github.com/aserto-dev/azm/cache"
	"github.com/aserto-dev/azm/model"
	dsc3 "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr3 "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
)

type SafeGetGraph struct {
	*dsr3.GetGraphRequest
}

func GetGraph(i *dsr3.GetGraphRequest) *SafeGetGraph {
	return &SafeGetGraph{i}
}

func (i *SafeGetGraph) Anchor() *dsc3.ObjectIdentifier {
	return &dsc3.ObjectIdentifier{
		ObjectType: i.AnchorType,
		ObjectId:   i.AnchorId,
	}
}

func (i *SafeGetGraph) Object() *dsc3.ObjectIdentifier {
	return &dsc3.ObjectIdentifier{
		ObjectType: i.ObjectType,
		ObjectId:   i.ObjectId,
	}
}

func (i *SafeGetGraph) Subject() *dsc3.ObjectIdentifier {
	return &dsc3.ObjectIdentifier{
		ObjectType: i.SubjectType,
		ObjectId:   i.SubjectId,
	}
}

func (i *SafeGetGraph) Validate(mc *cache.Cache) error {
	if i == nil || i.GetGraphRequest == nil {
		return derr.ErrInvalidRequest.Msg("get_graph")
	}

	// anchor must be defined, hence use an ObjectIdentifier.
	if err := ObjectIdentifier(i.Anchor()).Validate(mc); err != nil {
		return err
	}

	// Object can be optional, hence the use of an ObjectSelector.
	if err := ObjectSelector(i.Object()).Validate(mc); err != nil {
		return err
	}

	// Relation can be optional, hence the use of a RelationTypeSelector.
	if i.GetRelation() != "" {
		rr := &model.RelationRef{
			Object:   model.ObjectName(i.ObjectType),
			Relation: model.RelationName(i.Relation),
		}
		if err := RelationIdentifier(rr).Validate(AsRelation, mc); err != nil {
			return err
		}
	}

	// Subject can be option, hence the use of an ObjectSelector.
	if err := ObjectSelector(i.Subject()).Validate(mc); err != nil {
		return err
	}

	// either Object or Subject must be equal to the Anchor to indicate the directionality of the graph walk.
	// Anchor == Subject ==> subject->object (this was the default and only directionality before enabling bi-directionality)
	// Anchor == Object ==> object->subject
	if !ObjectIdentifier(i.Anchor()).Equal(i.Object()) &&
		!ObjectIdentifier(i.Anchor()).Equal(i.Subject()) {
		return derr.ErrGraphDirectionality
	}

	return nil
}
