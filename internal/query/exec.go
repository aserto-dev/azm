package query

import (
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

	ObjSet = ds.Set[model.ObjectID]
)

type Scope struct {
	OID model.ObjectID
	SID model.ObjectID
}

type PathSet = ds.Set[Scope]

type State interface {
	AddSet(ObjSet)
	ShortCircuit() bool
	Scopes() []Scope
	Result() ObjSet
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
	cache  Cache
}

func NewInterpreter(plan *Plan, getRels RelationReader, pool *mempool.RelationsPool) *Interpreter {
	return &Interpreter{
		plan:   plan,
		loader: newRelationLoader(getRels, pool),
		pool:   pool,
		cache:  Cache{},
	}
}

func (i *Interpreter) Run(req *dsr.CheckRequest) (ObjSet, error) {
	i.state = ds.NewStack[State](NewCompositeState(Union, 1, []Scope{{ObjectID(req.ObjectId), ObjectID(req.SubjectId)}}))
	if err := i.plan.Visit(i); err != nil {
		return nil, err
	}

	if i.state.Len() != 1 {
		return nil, errors.Wrap(ErrInterpreter, "unbalanced stack")
	}

	return i.state.Top().Result(), nil
}

func (i *Interpreter) OnLoad(expr *Load) error {
	state := i.state.Top()

	for _, scope := range state.Scopes() {
		if state.ShortCircuit() {
			return nil
		}

		if expr.Modifier.Has(SubjectWildcard) {
			scope.SID = "*"
		}
		if expr.Modifier.Has(ObjectWildcard) {
			scope.OID = "*"
		}

		result, err := i.loadSet(&Relation{RelationType: expr.RelationType, Scope: scope})
		if err != nil {
			return err
		}

		state.AddSet(result)
	}

	return nil
}

func (i *Interpreter) OnPipeStart(pipe *Pipe) (StepOption, error) {
	state := i.state.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}
	i.state.Push(NewChainState(state.Scopes()))
	return StepInto, nil
}

func (i *Interpreter) OnPipeEnd(_ *Pipe) {
	i.rollupResult()
}

func (i *Interpreter) OnCallStart(call *Call) (StepOption, error) {
	state := i.state.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}

	i.state.Push(NewCallState(call.Signature, i.state.Top().Scopes(), i.cache))

	return StepInto, nil
}

func (i *Interpreter) OnCallEnd(_ *Call) {
	i.rollupResult()
}

func (i *Interpreter) OnCompositeStart(expr *Composite) (StepOption, error) {
	state := i.state.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}

	i.state.Push(NewCompositeState(expr.Operator, len(expr.Operands), i.state.Top().Scopes()))
	return StepInto, nil
}

func (i *Interpreter) OnCompositeEnd(_ *Composite) {
	i.rollupResult()
}

func (i *Interpreter) rollupResult() {
	if i.state.Len() > 1 {
		state := i.state.Pop()
		i.state.Top().AddSet(state.Result())
	}
}

func (i *Interpreter) loadSet(rel *Relation) (ObjSet, error) {
	if result, ok := i.cache.LookupSet(rel); ok {
		if result == nil {
			result = ds.NewSet[model.ObjectID]()
		}
		return result, nil
	}

	var rid dsc.RelationIdentifier
	rel.Identifier(&rid)

	relsPtr := i.pool.GetSlice()
	if err := i.loader(&rid, relsPtr); err != nil {
		return nil, err
	}

	resultSet := ds.SetFromSlice(*relsPtr, func(rid *dsc.RelationIdentifier) model.ObjectID {
		return model.ObjectID(rid.SubjectId)
	})

	i.pool.PutSlice(relsPtr)

	i.cache.StoreSet(rel, resultSet)

	return resultSet, nil
}
