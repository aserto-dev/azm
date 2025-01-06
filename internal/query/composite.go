package query

type CompositeContext struct {
	op        Operator
	size      int
	remaining int
	hasResult bool
	scopes    []Scope
	result    ObjSet
}

func newCompositeContext(op Operator, size int, scopes []Scope, result ObjSet) *CompositeContext {
	return &CompositeContext{
		op:        op,
		size:      size,
		remaining: size,
		result:    result,
		scopes:    scopes,
	}
}

func (s *CompositeContext) AddSet(result ObjSet) {
	s.remaining--

	switch s.op {
	case Union:
		s.result.Add(result.Elements())
		if !s.result.IsEmpty() || s.remaining == 0 {
			// either we found a hit or exhausted all options.
			s.hasResult = true
		}
	case Intersection:
		if s.result.IsEmpty() {
			s.result.Union(result)
		} else {
			s.result.Intersect(result)
		}
		if s.result.IsEmpty() || s.remaining == 0 {
			// we either found a miss or exhausted all options.
			s.hasResult = true
		}
	case Difference:
		isFirst := s.remaining+1 == s.size
		if isFirst {
			s.result.Union(result)
		} else {
			s.result.Difference(result)
		}

		if s.result.IsEmpty() || s.remaining == 0 {
			// we either found a miss or exhausted all options.
			s.hasResult = true
		}
	}
}

func (s *CompositeContext) ShortCircuit() bool {
	return s.hasResult
}

func (s *CompositeContext) Scopes() []Scope {
	return s.scopes
}

func (s *CompositeContext) Result() ObjSet {
	return s.result
}
