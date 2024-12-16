package query

import (
	"slices"

	"github.com/hashicorp/go-set"

	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"

	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
)

type MessagePool[T any] interface {
	Get() T
	Put(T)
}

type (
	ObjectID     = model.ObjectID
	Relations    = []*dsc.RelationIdentifier
	RelationPool = MessagePool[*dsc.RelationIdentifier]

	// RelationReader retrieves relations that match the given filter.
	RelationReader func(*dsc.RelationIdentifier, RelationPool, *Relations) error
)

type step struct {
	cond Expression
	oid  model.ObjectID
	sid  model.ObjectID
}

func Exec(
	req *dsr.CheckRequest,
	m *model.Model,
	plan *Plan,
	getRels RelationReader,
	pool *mempool.RelationsPool,
) (bool, error) {
	backlog := ds.NewStack[step]()
	backlog.Push(step{plan.Expression, ObjectID(req.ObjectId), ObjectID(req.SubjectId)})

	stack := ds.NewStack[*marker]()
	stack.Push(newMarker(Union, 1))

	for !backlog.IsEmpty() {
		cur := backlog.Pop()

		switch cond := cur.cond.(type) {
		case Composite:
			stack.Pop()
			stack.Push(newMarker(cond.Operator, len(cond.Operands)))

			// push operands in reverse order so they pop out in the right order.
			for _, op := range slices.Backward(cond.Operands) {
				backlog.Push(step{op, cur.oid, cur.sid})
			}
		// case Call:
		// fun := plan.Functions[cond.Signature]
		//
		// rid := &dsc.RelationIdentifier{
		// 	ObjectType: cond.Param.OT.String(),
		// }
		// backlog.Push(step{fun,

		case Set:
			curMarker := stack.Top()
			if curMarker.HasResult() {
				// Short circuit. We already have a result.
				continue
			}

			rid := &dsc.RelationIdentifier{
				ObjectType:  cond.OT.String(),
				ObjectId:    cur.oid.String(),
				Relation:    cond.RT.String(),
				SubjectType: cond.ST.String(),
				SubjectId:   cur.sid.String(),
			}

			relsPtr := pool.GetSlice()
			if err := getRels(rid, pool, relsPtr); err != nil {
				return false, err
			}

			resultSet := set.FromFunc(*relsPtr, func(rid *dsc.RelationIdentifier) model.ObjectID {
				return model.ObjectID(rid.SubjectId)
			})

			curMarker.CompleteStep(resultSet)
			pool.PutSlice(relsPtr)

			for curMarker.IsDone() && stack.Len() > 1 {
				// The current condition is done. Roll up the result.
				stack.Pop()
				prevMarker := stack.Top()
				prevMarker.CompleteStep(curMarker.Result())
				curMarker = prevMarker
			}

		}
	}

	if stack.Len() != 1 {
		panic("unbalanced stack")
	}

	return !stack.Top().Result().Empty(), nil
}

type decision int

const (
	dPending decision = iota
	dFalse
	dTrue
)

type ObjSet = set.Set[model.ObjectID]

type marker struct {
	op        Operator
	size      int
	remaining int
	hasResult bool
	result    *ObjSet
}

func newMarker(op Operator, size int) *marker {
	return &marker{
		op:        op,
		size:      size,
		remaining: size,
		result:    set.New[model.ObjectID](1),
	}
}

func (m *marker) HasResult() bool {
	return m.hasResult
}

func (m *marker) Result() *ObjSet {
	return m.result
}

func (m *marker) IsDone() bool {
	return m.remaining == 0
}

func (m *marker) CompleteStep(resultSet *ObjSet) {
	m.remaining--

	switch m.op {
	case Union:
		m.result = m.result.Union(resultSet)
		if !m.result.Empty() || m.remaining == 0 {
			// either we found a hit or exhausted all options.
			m.hasResult = true
		}
	case Intersection:
		if m.result.Empty() {
			m.result = resultSet
		} else {
			m.result = m.result.Intersect(resultSet)
		}
		if m.result.Empty() || m.remaining == 0 {
			// we either found a miss or exhausted all options.
			m.hasResult = true
		}
	case Difference:
		isFirst := m.remaining+1 == m.size
		if isFirst {
			m.result = resultSet
		} else {
			m.result = m.result.Difference(resultSet)
		}

		if m.result.Empty() || m.remaining == 0 {
			// we either found a miss or exhausted all options.
			m.hasResult = true
		}
	}
}
