package query

import (
	"github.com/samber/lo"
)

type CallContext struct {
	signature *RelationType
	scopes    []Scope
	result    ObjSet
	cache     Cache
}

func newCallContext(sig *RelationType, scopes []Scope, result ObjSet, cache Cache) *CallContext {
	scopes = lo.Filter(scopes, func(p Scope, _ int) bool {
		key := Relation{
			RelationType: sig,
			Scope:        &p,
		}
		if res, ok := cache.LookupCall(&key); ok {
			if res != nil {
				result.Union(*res)
			}
			return false
		}

		cache.StoreCall(&key, nil)
		return true
	})

	return &CallContext{
		signature: sig,
		scopes:    scopes,
		result:    result,
		cache:     cache,
	}
}

func (s *CallContext) AddSet(result ObjSet) {
	scope := s.scopes[0]
	s.scopes = s.scopes[1:]

	key := Relation{
		RelationType: s.signature,
		Scope:        &scope,
	}
	s.cache.StoreCall(&key, &result)
	s.result.Union(result)
}

func (s *CallContext) ShortCircuit() bool {
	return !s.result.IsEmpty()
}

func (s *CallContext) Scopes() []Scope {
	return s.scopes
}

func (s *CallContext) Result() ObjSet {
	return s.result
}
