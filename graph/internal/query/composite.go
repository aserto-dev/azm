package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/hashicorp/go-set"
)

type CompositeState struct {
	op        Operator
	size      int
	remaining int
	hasResult bool
	scopes    []scope
	result    *ObjSet
}

func NewCompositeState(op Operator, size int, scopes []scope) *CompositeState {
	return &CompositeState{
		op:        op,
		size:      size,
		remaining: size,
		result:    set.New[model.ObjectID](1),
		scopes:    scopes,
	}
}

func (m *CompositeState) AddResult(result *ObjSet) {
	m.remaining--

	switch m.op {
	case Union:
		m.result = m.result.Union(result)
		if !m.result.Empty() || m.remaining == 0 {
			// either we found a hit or exhausted all options.
			m.hasResult = true
		}
	case Intersection:
		if m.result.Empty() {
			m.result = result
		} else {
			m.result = m.result.Intersect(result)
		}
		if m.result.Empty() || m.remaining == 0 {
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

		if m.result.Empty() || m.remaining == 0 {
			// we either found a miss or exhausted all options.
			m.hasResult = true
		}
	}
}

func (m *CompositeState) ShortCircuit() bool {
	return m.hasResult
}

func (m *CompositeState) Scopes() []scope {
	return m.scopes
}

func (m *CompositeState) Result() *ObjSet {
	return m.result
}
