package query

import (
	"slices"

	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type CallState struct {
	signature *Set
	paths     []Path
	calls     PathSet
	pending   []Path
	result    ObjSet
	cache     Cache
}

func NewCallState(sig *Set, paths []Path, cache Cache) *CallState {
	return &CallState{
		signature: sig,
		paths:     paths,
		calls:     NewSet[Path](),
		result:    NewSet[model.ObjectID](),
		cache:     cache,
	}
}

func (s *CallState) AddResult(result ObjSet) {
	if len(s.paths) > 0 {
		basePath := s.paths[0]
		s.paths = s.paths[1:]

		for oid := range result.Iter() {
			callPath := Path{OID: oid, SID: basePath.SID}
			key := Relation{
				Set:  *s.signature,
				Path: callPath,
			}
			if res, ok := s.cache.LookupCall(&key); ok {
				if res != nil {
					s.result = s.result.Union(res)
				}
			} else {
				s.cache.StoreCall(&key, nil)
				s.calls.Add(Path{oid, basePath.SID})
			}
		}

		return
	}

	if len(s.pending) == 0 {
		panic("unexpected call to AddResult. no pending calls")
	}
	path := s.pending[0]
	s.pending = s.pending[1:]

	key := Relation{
		Set:  *s.signature,
		Path: path,
	}
	s.cache.StoreCall(&key, result)
	s.result = s.result.Union(result)
}

func (s *CallState) ShortCircuit() bool {
	expanded := len(s.paths) == 0
	pending := len(s.pending) > 0 || !s.calls.IsEmpty()
	return expanded && (!pending || !s.result.IsEmpty())
}

func (s *CallState) Paths() []Path {
	if len(s.paths) > 0 {
		return lo.Map(s.paths, func(s Path, _ int) Path {
			return Path{OID: s.OID}
		})
	}

	if s.pending == nil {
		s.pending = slices.SortedFunc(slices.Values(s.calls.ToSlice()), func(lhs, rhs Path) int {
			if lhs.OID < rhs.OID {
				return -1
			}
			if lhs.OID > rhs.OID {
				return 1
			}
			return 0
		})
	}

	return s.pending
}

func (s *CallState) Result() ObjSet {
	return s.result
}

func (s *CallState) emptyResult() bool {
	return s.result == nil || s.result.IsEmpty()
}
