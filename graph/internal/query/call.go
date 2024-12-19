package query

type CallState struct {
	Func Expression
	Args *ObjSet
}

func NewCallState(f Expression) *CallState {
	return &CallState{Func: f}
}

func (m *CallState) ShortCircuit() bool {
	return false
}

func (m *CallState) Result() *ObjSet {
	return nil
}

func (m *CallState) IsDone() bool {
	return false
}

func (m *CallState) AddResult(resultSet *ObjSet) {
	m.Args = resultSet
}
