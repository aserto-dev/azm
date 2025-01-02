package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type CallState struct {
	signature *RelationType
	scopes    []Scope
	result    ObjSet
	cache     Cache
}

func NewCallState(sig *RelationType, scopes []Scope, cache Cache) *CallState {
	result := NewSet[model.ObjectID]()

	scopes = lo.Filter(scopes, func(p Scope, _ int) bool {
		key := Relation{
			RelationType: sig,
			Scope:        p,
		}
		if res, ok := cache.LookupCall(&key); ok {
			if res != nil {
				result = result.Union(res)
			}
			return false
		}

		cache.StoreCall(&key, nil)
		return true
	})

	return &CallState{
		signature: sig,
		scopes:    scopes,
		result:    result,
		cache:     cache,
	}
}

func (s *CallState) AddSet(result ObjSet) {
	path := s.scopes[0]
	s.scopes = s.scopes[1:]

	key := Relation{
		RelationType: s.signature,
		Scope:        path,
	}
	s.cache.StoreCall(&key, result)
	s.result = s.result.Union(result)
}

func (s *CallState) ShortCircuit() bool {
	return !s.result.IsEmpty()
}

func (s *CallState) Scopes() []Scope {
	return s.scopes
}

func (s *CallState) Result() ObjSet {
	return s.result
}
