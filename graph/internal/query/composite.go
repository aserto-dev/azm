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

func (m *CompositeState) AddResult(result ObjSet) {
	m.remaining--

	switch m.op {
	case Union:
		m.result = m.result.Union(result)
		if !m.result.IsEmpty() || m.remaining == 0 {
			// either we found a hit or exhausted all options.
			m.hasResult = true
		}
	case Intersection:
		if m.result.IsEmpty() {
			m.result = result
		} else {
			m.result = m.result.Intersect(result)
		}
		if m.result.IsEmpty() || m.remaining == 0 {
			// we either found a miss or exhausted all options.
			m.hasResult = true
		}
	case Difference:
		isFirst := m.remaining+1 == m.size
		if isFirst {
			m.result = result
		} else {
			m.result = m.result.Difference(result)
		}

		if m.result.IsEmpty() || m.remaining == 0 {
			// we either found a miss or exhausted all options.
			m.hasResult = true
		}
	}
}

func (m *CompositeState) ShortCircuit() bool {
	return m.hasResult
}

func (m *CompositeState) Paths() []Path {
	return m.paths
}

func (m *CompositeState) Result() ObjSet {
	return m.result
}
