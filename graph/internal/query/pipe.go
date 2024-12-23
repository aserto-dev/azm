package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type ChainState struct {
	paths     []Path
	expansion PathSet
	result    ObjSet
}

func NewChainState(paths []Path) *ChainState {
	return &ChainState{paths: paths, expansion: NewSet[Path](), result: NewSet[model.ObjectID]()}
}

func (s *ChainState) AddSet(result ObjSet) {
	if len(s.paths) > 0 {
		path := s.paths[0]
		s.paths = s.paths[1:]

		for oid := range result.Iter() {
			s.expansion.Add(Path{oid, path.SID})
		}

		return
	}

	s.result = s.result.Union(result)
}

func (s *ChainState) ShortCircuit() bool {
	return len(s.paths) == 0 && s.expansion.IsEmpty()
}

func (s *ChainState) Paths() []Path {
	if len(s.paths) > 0 {
		return lo.Map(s.paths, func(p Path, _ int) Path {
			return Path{OID: p.OID}
		})
	}

	return s.expansion.ToSlice()
}

func (s *ChainState) Result() ObjSet { return s.result }
