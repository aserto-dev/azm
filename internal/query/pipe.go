package query

import (
	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type ScopeSet = ds.Set[Scope]

type PipeContext struct {
	scopes    []Scope
	expansion ScopeSet
	result    ObjSet
}

func newPipeContext(scopes []Scope, result ObjSet) *PipeContext {
	return &PipeContext{scopes: scopes, expansion: ds.NewSet[Scope](len(scopes)), result: result}
}

func (s *PipeContext) AddSet(result ObjSet) {
	if len(s.scopes) > 0 {
		path := s.scopes[0]
		s.scopes = s.scopes[1:]

		s.expansion.Add(ds.TransformIter(result.Elements(), func(oid model.ObjectID) Scope {
			return Scope{oid, path.SID}
		}))

		return
	}

	s.result.Union(result)
}

func (s *PipeContext) ShortCircuit() bool {
	return len(s.scopes) == 0 && s.expansion.IsEmpty()
}

func (s *PipeContext) Scopes() []Scope {
	if len(s.scopes) > 0 {
		return lo.Map(s.scopes, func(p Scope, _ int) Scope {
			return Scope{OID: p.OID}
		})
	}

	return s.expansion.ToSlice()
}

func (s *PipeContext) Result() ObjSet {
	return s.result
}
