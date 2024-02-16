package check

import (
	"fmt"
	"strings"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/structpb"
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
	str := fmt.Sprintf("%s:%s#%s@%s:%s", p.ot, displayID(p.oid), p.rel, p.st, displayID(p.sid))
	if p.srel != "" {
		str += fmt.Sprintf("#%s", p.srel)
	}
	return str
}

func displayID(id ObjectID) string {
	if id == "" {
		return "?"
	}
	return id.String()
}

// converts a relation to a dsc.Relation.
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

func (p *relation) Object() *object {
	return &object{
		Type: p.ot,
		ID:   p.oid,
	}
}

func (p *relation) Subject() *object {
	return &object{
		Type: p.st,
		ID:   p.sid,
	}
}

type relations []*relation

type searchPath relations

type object struct {
	Type model.ObjectName
	ID   ObjectID
}

// The results of a search is a map where the key is a matching relations
// and the value is a list of paths that connect the search object and subject.
type searchResults map[object][]searchPath

// Objects returns the objects from the search results.
func (r searchResults) Objects() []*dsc.ObjectIdentifier {
	return lo.MapToSlice(r, func(o object, _ []searchPath) *dsc.ObjectIdentifier {
		return &dsc.ObjectIdentifier{
			ObjectType: o.Type.String(),
			ObjectId:   o.ID.String(),
		}
	})
}

// Subjects returns the subjects from the search results.
func (r searchResults) Subjects() []*dsc.ObjectIdentifier {
	return lo.MapToSlice(r, func(o object, _ []searchPath) *dsc.ObjectIdentifier {
		return &dsc.ObjectIdentifier{
			ObjectType: o.Type.String(),
			ObjectId:   o.ID.String(),
		}
	})
}

func (r searchResults) Explain() *structpb.Struct {
	explanation := lo.MapEntries(r, func(obj object, paths []searchPath) (string, any) {
		key := fmt.Sprintf("%s:%s", obj.Type, obj.ID)

		val := lo.Map(paths, func(path searchPath, _ int) any {
			return lo.Map(path, func(rel *relation, _ int) any {
				return rel.String()
			})
		})

		return key, val
	})

	res, err := structpb.NewStruct(explanation)
	if err != nil {
		panic(err)
	}

	return res
}

type searchStatus int

func (s searchStatus) String() string {
	switch s {
	case searchStatusUnknown:
		return statusUnknown
	case searchStatusPending:
		return statusPending
	case searchStatusComplete:
		return statusComplete
	default:
		return fmt.Sprintf("invalid: %d", s)
	}
}

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

type searchCall struct {
	*relation
	status searchStatus
}

type searchMemo struct {
	visited map[relation]searchResults
	history []*searchCall
}

func newSearchMemo(trace bool) *searchMemo {
	return &searchMemo{
		visited: map[relation]searchResults{},
		history: lo.Ternary(trace, []*searchCall{}, nil),
	}
}

func (m *searchMemo) MarkVisited(params *relation) searchStatus {
	results, ok := m.visited[*params]
	switch {
	case !ok:
		m.visited[*params] = nil
		m.trace(params, searchStatusPending)
		return searchStatusUnknown
	case results == nil:
		return searchStatusPending
	default:
		return searchStatusComplete
	}
}

func (m *searchMemo) MarkComplete(params *relation, results searchResults) {
	m.visited[*params] = results
	m.trace(params, searchStatusComplete)

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

func (m *searchMemo) Trace() []string {
	if m.history == nil {
		return []string{}
	}

	callstack := []string{}

	return lo.Map(m.history, func(c *searchCall, _ int) string {
		call := c.String()
		result := c.status.String()

		if len(callstack) > 0 && callstack[len(callstack)-1] == call && c.status != searchStatusPending {
			callstack = callstack[:len(callstack)-1]
		}

		s := fmt.Sprintf("%s%s = %s", strings.Repeat("  ", len(callstack)), call, result)

		if c.status == searchStatusPending {
			callstack = append(callstack, call)
		}

		return s
	})
}

func (m *searchMemo) trace(params *relation, status searchStatus) {
	if m.history != nil {
		m.history = append(m.history, &searchCall{params, status})
	}
}
