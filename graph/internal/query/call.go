package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/hashicorp/go-set"
	"github.com/samber/lo"
)

type CallState struct {
	function Expression
	params   *ObjSet
	scopes   []scope
	result   *ObjSet
}

func NewCallState(f Expression, scopes []scope) *CallState {
	return &CallState{function: f, scopes: scopes, params: set.New[model.ObjectID](1), result: set.New[model.ObjectID](1)}
}

func (m *CallState) AddResult(resultSet *ObjSet) {
	if m.params == nil {
		m.params = m.params.Union(resultSet)
		return
	}

	m.result = m.result.Union(resultSet)
}

func (m *CallState) ShortCircuit() bool {
	return !m.hasParams() || m.emptyResult()
}

func (m *CallState) Scopes() []scope {
	if m.params == nil {
		return lo.Map(m.scopes, func(s scope, _ int) scope {
			return scope{ObjectID: s.ObjectID}
		})
	}

	scopes := make([]scope, 0, m.params.Size()*len(m.scopes))
	for _, oid := range m.params.Slice() {
		for _, s := range m.scopes {
			scopes = append(scopes, scope{ObjectID: oid, SubjectID: s.SubjectID})
		}
	}

	return scopes
}

func (m *CallState) Result() *ObjSet {
	return m.result
}

func (m *CallState) hasParams() bool {
	return m.params != nil && !m.params.Empty()
}

func (m *CallState) emptyResult() bool {
	return m.result == nil || m.result.Empty()
}
