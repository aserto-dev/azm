package check

import (
	"fmt"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/samber/lo"
)

type ObjectID = model.ObjectID

// RelationReader retrieves relations that match the given filter.
type RelationReader func(*dsc.Relation) ([]*dsc.Relation, error)

// relation is a comparable representation of a relation. It can be used as a map key.
type relation struct {
	ot   model.ObjectName
	oid  ObjectID
	rel  model.RelationName
	st   model.ObjectName
	sid  ObjectID
	srel model.RelationName
}

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
	str := fmt.Sprintf("%s:%s#%s@%s:%s", p.ot, p.oid, p.rel, p.st, p.sid)
	if p.srel != "" {
		str += fmt.Sprintf("#%s", p.srel)
	}
	return str
}

// convers a relation to a dsc.Relation.
func (p *relation) AsRelation() *dsc.Relation {
	return &dsc.Relation{
		ObjectType:      p.ot.String(),
		ObjectId:        p.oid.String(),
		Relation:        p.rel.String(),
		SubjectType:     p.st.String(),
		SubjectId:       p.sid.String(),
		SubjectRelation: p.srel.String(),
	}
}

type relations []*relation

type searchPath relations

// The results of a search is a map where the key is a matching relations
// and the value is a list of paths that connect the search object and subject.
type searchResults map[relation][]searchPath

// Objects returns the objects from the search results.
func (r searchResults) Objects() []*dsc.ObjectIdentifier {
	return lo.MapToSlice(r, func(p relation, _ []searchPath) *dsc.ObjectIdentifier {
		return &dsc.ObjectIdentifier{
			ObjectType: p.ot.String(),
			ObjectId:   p.oid.String(),
		}
	})
}

// Subjects returns the subjects from the search results.
func (r searchResults) Subjects() []*dsc.ObjectIdentifier {
	return lo.MapToSlice(r, func(p relation, _ []searchPath) *dsc.ObjectIdentifier {
		return &dsc.ObjectIdentifier{
			ObjectType: p.st.String(),
			ObjectId:   p.sid.String(),
		}
	})
}

type searchStatus int

const (
	searchStatusUnknown searchStatus = iota
	searchStatusPending
	searchStatusComplete
)

type graphSearch struct {
	m       *model.Model
	params  *relation
	getRels RelationReader

	memo    *searchMemo
	explain bool
}

func (s *graphSearch) validate() error {
	o := s.m.Objects[s.params.ot]
	if o == nil {
		return derr.ErrObjectTypeNotFound.Msg(s.params.ot.String())
	}

	if !o.HasRelOrPerm(s.params.rel) {
		return derr.ErrRelationNotFound.Msg(s.params.rel.String())
	}

	if _, ok := s.m.Objects[s.params.st]; !ok {
		return derr.ErrObjectTypeNotFound.Msg(s.params.st.String())
	}

	return nil
}

func (s *graphSearch) Explain() searchResults {
	return s.memo.visited[*s.params]
}

type searchMemo struct {
	visited map[relation]searchResults
	history relations
}

func newSearchMemo(trace bool) *searchMemo {
	return &searchMemo{
		visited: map[relation]searchResults{},
		history: lo.Ternary(trace, relations{}, nil),
	}
}

func (m *searchMemo) MarkVisited(params *relation) searchStatus {
	results, ok := m.visited[*params]
	switch {
	case !ok:
		m.visited[*params] = nil
		if m.history != nil {
			m.history = append(m.history, params)
		}
		return searchStatusUnknown
	case results == nil:
		return searchStatusPending
	default:
		return searchStatusComplete
	}
}

func (m *searchMemo) MarkComplete(params *relation, results searchResults) {
	m.visited[*params] = results
}

func (m *searchMemo) Status(params *relation) searchStatus {
	results, ok := m.visited[*params]
	switch {
	case !ok:
		return searchStatusUnknown
	case results == nil:
		return searchStatusPending
	default:
		return searchStatusComplete
	}
}
