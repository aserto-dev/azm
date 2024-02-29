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

// converts a relation to a dsc.Relation.
func (p *relation) toProto() *dsc.Relation {
	return &dsc.Relation{
		ObjectType:      p.ot.String(),
		ObjectId:        p.oid.String(),
		Relation:        p.rel.String(),
		SubjectType:     p.st.String(),
		SubjectId:       p.sid.String(),
		SubjectRelation: p.srel.String(),
	}
}

func (p *relation) String() string {
	str := fmt.Sprintf("%s:%s#%s@%s:%s", p.ot, displayID(p.oid), p.rel, p.st, displayID(p.sid))
	if p.srel != "" {
		str += fmt.Sprintf("#%s", p.srel)
	}
	return str
}

func (p *relation) object() *object {
	return &object{
		Type: p.ot,
		ID:   p.oid,
	}
}

func (p *relation) subject() *object {
	return &object{
		Type: p.st,
		ID:   p.sid,
	}
}

// func (p *relation) clone(overrides ...fieldOverride) *relation {
//     clone := *p

//     for _, override := range overrides {
//         override(&clone)
//     }

//     return &clone
// }

// type fieldOverride func(*relation)

// func withObjectType(ot model.ObjectName) fieldOverride {
//     return func(p *relation) {
//         p.ot = ot
//     }
// }

// func withObjectID(oid ObjectID) fieldOverride {
//     return func(p *relation) {
//         p.oid = oid
//     }
// }

// func withRelation(rel model.RelationName) fieldOverride {
//     return func(p *relation) {
//         p.rel = rel
//     }
// }

// func withSubjectType(st model.ObjectName) fieldOverride {
//     return func(p *relation) {
//         p.st = st
//     }
// }

// func withSubjectID(sid ObjectID) fieldOverride {
//     return func(p *relation) {
//         p.sid = sid
//     }
// }

// func withSubjectRelation(srel model.RelationName) fieldOverride {
//     return func(p *relation) {
//         p.srel = srel
//     }
// }

func displayID(id ObjectID) string {
	if id == "" {
		return "?"
	}
	return id.String()
}
