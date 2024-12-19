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
	plan *Plan,
	getRels RelationReader,
	pool *mempool.RelationsPool,
) (bool, error) {
	backlog := ds.NewStack(step{plan.Expression, ObjectID(req.ObjectId), ObjectID(req.SubjectId)})

	stack := ds.NewStack[State]()
	stack.Push(NewCompositeState(Union, 1))

	for !backlog.IsEmpty() {
		cur := backlog.Pop()

		switch cond := cur.cond.(type) {
		case Composite:
			stack.Pop()
			stack.Push(NewCompositeState(cond.Operator, len(cond.Operands)))

			// push operands in reverse order so they pop out in the right order.
			for _, op := range slices.Backward(cond.Operands) {
				backlog.Push(step{op, cur.oid, cur.sid})
			}
		case Call:
			fun := plan.Functions[cond.Signature]
			stack.Push(NewCallState(fun))
			backlog.Push(step{cond.Param, cur.oid, ""})
		case Set:
			curMarker := stack.Top()
			if curMarker.ShortCircuit() {
				// Short circuit. We already have a result.
				continue
			}

			rid := &dsc.RelationIdentifier{
				ObjectType:      cond.OT.String(),
				ObjectId:        cur.oid.String(),
				Relation:        cond.RT.String(),
				SubjectType:     cond.ST.String(),
				SubjectId:       cur.sid.String(),
				SubjectRelation: cond.SRT.String(),
			}

			relsPtr := pool.GetSlice()
			if err := getRels(rid, pool, relsPtr); err != nil {
				return false, err
			}

			resultSet := set.FromFunc(*relsPtr, func(rid *dsc.RelationIdentifier) model.ObjectID {
				return model.ObjectID(rid.SubjectId)
			})

			curMarker.AddResult(resultSet)
			pool.PutSlice(relsPtr)

			for curMarker.IsDone() && stack.Len() > 1 {
				// The current condition is done. Roll up the result.
				stack.Pop()
				prevMarker := stack.Top().(*CompositeState)
				prevMarker.AddResult(curMarker.Result())
				curMarker = prevMarker
			}
		}
	}

	if stack.Len() != 1 {
		panic("unbalanced stack")
	}

	return !stack.Top().Result().Empty(), nil
}

type ObjSet = set.Set[model.ObjectID]

type IDs struct {
	ObjectID  model.ObjectID
	SubjectID model.ObjectID
}

type State interface {
	AddResult(*ObjSet)
	ShortCircuit() bool
	IsDone() bool
	Result() *ObjSet
}

type RelationLoader func(*dsc.RelationIdentifier, *Relations) error

func newRelationLoader(reader RelationReader, pool RelationPool) RelationLoader {
	return func(rid *dsc.RelationIdentifier, outRels *Relations) error {
		return reader(rid, pool, outRels)
	}
}

type Interpreter struct {
	plan   *Plan
	loader RelationLoader
	pool   *mempool.RelationsPool
	state  *ds.Stack[State]
	ids    *ds.Stack[IDs]
}

func NewInterpreter(plan *Plan, getRels RelationReader, pool *mempool.RelationsPool) *Interpreter {
	return &Interpreter{
		plan:   plan,
		loader: newRelationLoader(getRels, pool),
		pool:   pool,
	}
}

func (i *Interpreter) Run(req *dsr.CheckRequest) {
	i.state = ds.NewStack[State](NewCompositeState(Union, 1))
	i.ids = ds.NewStack(IDs{ObjectID(req.ObjectId), ObjectID(req.SubjectId)})
	i.plan.Visit(i)
}

func (i *Interpreter) VisitSet(expr Set) (bool, error) {
	ids := i.ids.Pop()

	rid := &dsc.RelationIdentifier{
		ObjectType:      expr.OT.String(),
		ObjectId:        ids.ObjectID.String(),
		Relation:        expr.RT.String(),
		SubjectType:     expr.ST.String(),
		SubjectId:       ids.SubjectID.String(),
		SubjectRelation: expr.SRT.String(),
	}

	relsPtr := i.pool.GetSlice()
	if err := i.loader(rid, relsPtr); err != nil {
		return false, err
	}

	resultSet := set.FromFunc(*relsPtr, func(rid *dsc.RelationIdentifier) model.ObjectID {
		return model.ObjectID(rid.SubjectId)
	})

	state := i.state.Top()
	state.AddResult(resultSet)

	return !state.ShortCircuit(), nil
}

func (i *Interpreter) VisitCall(call Call) (bool, error) {
	return true, nil
}

func (i *Interpreter) VisitComposite(expr Composite) (bool, error) {
	return true, nil
}
