package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

type PipeState struct {
	paths     []Path
	expansion PathSet
	result    ObjSet
}

func NewChainState(paths []Path) *PipeState {
	return &PipeState{paths: paths, expansion: NewSet[Path](), result: NewSet[model.ObjectID]()}
}

func (s *PipeState) AddSet(result ObjSet) {
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

func (s *PipeState) ShortCircuit() bool {
	return len(s.paths) == 0 && s.expansion.IsEmpty()
}

func (s *PipeState) Paths() []Path {
	if len(s.paths) > 0 {
		return lo.Map(s.paths, func(p Path, _ int) Path {
			return Path{OID: p.OID}
		})
	}

	return s.expansion.ToSlice()
}

func (s *PipeState) Result() ObjSet { return s.result }
