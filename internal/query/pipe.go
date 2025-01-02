package query

import (
	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type PipeState struct {
	scopes    []Scope
	expansion PathSet
	result    ObjSet
}

func NewChainState(scopes []Scope) *PipeState {
	return &PipeState{scopes: scopes, expansion: ds.NewSet[Scope](), result: ds.NewSet[model.ObjectID]()}
}

func (s *PipeState) AddSet(result ObjSet) {
	if len(s.scopes) > 0 {
		path := s.scopes[0]
		s.scopes = s.scopes[1:]

		for oid := range result.Elements() {
			s.expansion.Add(Scope{oid, path.SID})
		}

		return
	}

	s.result = s.result.Union(result)
}

func (s *PipeState) ShortCircuit() bool {
	return len(s.scopes) == 0 && s.expansion.IsEmpty()
}

func (s *PipeState) Scopes() []Scope {
	if len(s.scopes) > 0 {
		return lo.Map(s.scopes, func(p Scope, _ int) Scope {
			return Scope{OID: p.OID}
		})
	}

	return s.expansion.ToSlice()
}

func (s *PipeState) Result() ObjSet { return s.result }
