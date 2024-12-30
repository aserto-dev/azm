package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type CallState struct {
	signature *Load
	paths     []Path
	result    ObjSet
	cache     Cache
}

func NewCallState(sig *Load, paths []Path, cache Cache) *CallState {
	result := NewSet[model.ObjectID]()

	paths = lo.Filter(paths, func(p Path, _ int) bool {
		key := Relation{
			Load: *sig,
			Path: p,
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
		paths:     paths,
		result:    result,
		cache:     cache,
	}
}

func (s *CallState) AddSet(result ObjSet) {
	path := s.paths[0]
	s.paths = s.paths[1:]

	key := Relation{
		Load: *s.signature,
		Path: path,
	}
	s.cache.StoreCall(&key, result)
	s.result = s.result.Union(result)
}

func (s *CallState) ShortCircuit() bool {
	return !s.result.IsEmpty()
}

func (s *CallState) Paths() []Path {
	return s.paths
}

func (s *CallState) Result() ObjSet {
	return s.result
}
