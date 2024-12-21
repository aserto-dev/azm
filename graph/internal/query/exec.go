package query

import (
	"github.com/hashicorp/go-set"
	"github.com/pkg/errors"

	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"

	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
)

var ErrInterpreter = errors.New("query interpreter error")

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

	ObjSet = set.Set[model.ObjectID]
)

type scope struct {
	ObjectID  model.ObjectID
	SubjectID model.ObjectID
}

type State interface {
	AddResult(*ObjSet)
	ShortCircuit() bool
	Scopes() []scope
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
}

func NewInterpreter(plan *Plan, getRels RelationReader, pool *mempool.RelationsPool) *Interpreter {
	return &Interpreter{
		plan:   plan,
		loader: newRelationLoader(getRels, pool),
		pool:   pool,
	}
}

func (i *Interpreter) Run(req *dsr.CheckRequest) (*ObjSet, error) {
	i.state = ds.NewStack[State](NewCompositeState(Union, 1, []scope{{ObjectID(req.ObjectId), ObjectID(req.SubjectId)}}))
	if err := i.plan.Visit(i); err != nil {
		return nil, err
	}

	if i.state.Len() != 1 {
		return nil, errors.Wrap(ErrInterpreter, "unbalanced stack")
	}

	return i.state.Top().Result(), nil
}

func (i *Interpreter) OnSet(expr *Set) error {
	state := i.state.Top()

	for _, ids := range state.Scopes() {
		if state.ShortCircuit() {
			return nil
		}

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
			return err
		}

		resultSet := set.FromFunc(*relsPtr, func(rid *dsc.RelationIdentifier) model.ObjectID {
			return model.ObjectID(rid.SubjectId)
		})

		i.pool.PutSlice(relsPtr)

		state.AddResult(resultSet)
	}

	return nil
}

func (i *Interpreter) OnCallStart(call *Call) (VisitOption, error) {
	state := i.state.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}

	f := i.plan.Functions[*call.Signature]
	i.state.Push(NewCallState(f, i.state.Top().Scopes()))

	return StepInto, nil
}

func (i *Interpreter) OnCallEnd(_ *Call) {
	i.rollupResult()
}

func (i *Interpreter) OnCompositeStart(expr *Composite) error {
	i.state.Push(NewCompositeState(expr.Operator, len(expr.Operands), i.state.Top().Scopes()))
	return nil
}

func (i *Interpreter) OnCompositeEnd(_ *Composite) {
	i.rollupResult()
}

func (i *Interpreter) rollupResult() {
	if i.state.Len() > 1 {
		state := i.state.Pop()
		i.state.Top().AddResult(state.Result())
	}
}
