package query

import (
	"slices"

	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
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
	cond Plan
	oid  model.ObjectID
	sid  model.ObjectID
}

func Exec(
	req *dsr.CheckRequest,
	m *model.Model,
	plan Plan,
	getRels RelationReader,
	pool *mempool.RelationsPool,
) (bool, error) {
	backlog := ds.NewStack[step]()
	backlog.Push(step{plan, ObjectID(req.ObjectId), ObjectID(req.SubjectId)})

	state := ds.NewStack[marker]()
	state.Push(marker{op: Union, len: 1})

	for !backlog.IsEmpty() {
		cur := backlog.Pop()

		switch cond := cur.cond.(type) {
		case Single:
			curMarker := state.Top()
			if curMarker.hasResult() {
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

			curMarker.completeStep(len(*relsPtr) > 0)
			pool.PutSlice(relsPtr)

			for curMarker.isDone() && state.Len() > 1 {
				// The current condition is done. Roll up the result.
				state.Pop()
				prevMarker := state.Top()
				prevMarker.completeStep(curMarker.result == dTrue)
				curMarker = prevMarker
			}

		case Composite:
			state.Pop()
			state.Push(marker{op: cond.Operator, len: len(cond.Operands)})
			// push operands in reverse order so they pop out in the right order.
			for _, op := range slices.Backward(cond.Operands) {
				backlog.Push(step{op, cur.oid, cur.sid})
			}
		}
	}

	if state.Len() != 1 {
		panic("unbalanced stack")
	}

	return state.Top().result == dTrue, nil
}

type decision int

const (
	dPending decision = iota
	dFalse
	dTrue
)

type marker struct {
	op     Operator
	len    int
	result decision
}

func (m *marker) hasResult() bool {
	return m.result != dPending
}

func (m *marker) isDone() bool {
	return m.len == 0
}

func (m *marker) completeStep(result bool) {
	m.len--

	switch m.op {
	case Union:
		if result || m.len == 0 {
			// either we found a hit or exhausted all options.
			m.result = lo.Ternary(result, dTrue, dFalse)
		}
	case Intersection:
		if !result || m.len == 0 {
			// we either found a miss or exhausted all options.
			m.result = lo.Ternary(!result, dFalse, dTrue)
		}
	case Negation:
		if result || m.len == 0 {
			// we either found a miss or exhausted all options.
			m.result = lo.Ternary(result, dFalse, dTrue)
		}
	}
}
