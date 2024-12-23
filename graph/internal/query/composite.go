package query

import (
	"github.com/aserto-dev/azm/model"
)

type CompositeState struct {
	op        Operator
	size      int
	remaining int
	hasResult bool
	paths     []Path
	result    ObjSet
}

func NewCompositeState(op Operator, size int, paths []Path) *CompositeState {
	return &CompositeState{
		op:        op,
		size:      size,
		remaining: size,
		result:    NewSet[model.ObjectID](),
		paths:     paths,
	}
}

func (s *CompositeState) AddSet(result ObjSet) {
	s.remaining--

	switch s.op {
	case Union:
		s.result = s.result.Union(result)
		if !s.result.IsEmpty() || s.remaining == 0 {
			// either we found a hit or exhausted all options.
			s.hasResult = true
		}
	case Intersection:
		if s.result.IsEmpty() {
			s.result = result
		} else {
			s.result = s.result.Intersect(result)
		}
		if s.result.IsEmpty() || s.remaining == 0 {
			// we either found a miss or exhausted all options.
			s.hasResult = true
		}
	case Difference:
		isFirst := s.remaining+1 == s.size
		if isFirst {
			s.result = result
		} else {
			s.result = s.result.Difference(result)
		}

		if s.result.IsEmpty() || s.remaining == 0 {
			// we either found a miss or exhausted all options.
			s.hasResult = true
		}
	}
}

func (s *CompositeState) ShortCircuit() bool {
	return s.hasResult
}

func (s *CompositeState) Paths() []Path {
	return s.paths
}

func (s *CompositeState) Result() ObjSet {
	return s.result
}
