package graph

import (
	"fmt"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
)

// relation is a comparable representation of a relation. It can be used as a map key.
type relation struct {
	ot   model.ObjectName
	oid  ObjectID
	rel  model.RelationName
	st   model.ObjectName
	sid  ObjectID
	srel model.RelationName
}

type relations []*relation

// converts a dsc.Relation to a relation.
func relationFromProto(rel *dsc.Relation) *relation {
	return &relation{
		ot:   model.ObjectName(rel.ObjectType),
		oid:  ObjectID(rel.ObjectId),
		rel:  model.RelationName(rel.Relation),
		st:   model.ObjectName(rel.SubjectType),
		sid:  ObjectID(rel.SubjectId),
		srel: model.RelationName(rel.SubjectRelation),
	}
}

func (p *relation) String() string {
	str := fmt.Sprintf("%s:%s#%s@%s:%s", p.ot, displayID(p.oid), p.rel, p.st, displayID(p.sid))
	if p.srel != "" {
		str += fmt.Sprintf("#%s", p.srel)
	}
	return str
}

func (p *relation) subject() *object {
	return &object{
		Type: p.st,
		ID:   p.sid,
	}
}

func displayID(id ObjectID) string {
	if id == "" {
		return "?"
	}
	return id.String()
}
