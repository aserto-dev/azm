package query

import (
	mapset "github.com/deckarep/golang-set/v2"
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

	ObjSet = mapset.Set[model.ObjectID]
)

type Path struct {
	OID model.ObjectID
	SID model.ObjectID
}

type PathSet = mapset.Set[Path]

type State interface {
	AddResult(ObjSet)
	ShortCircuit() bool
	Paths() []Path
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
	i.state = ds.NewStack[State](NewCompositeState(Union, 1, []Path{{ObjectID(req.ObjectId), ObjectID(req.SubjectId)}}))
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

	for _, path := range state.Paths() {
		if state.ShortCircuit() {
			return nil
		}

		result, err := i.loadSet(&Relation{Set: *expr, Path: path})
		if err != nil {
			return err
		}

		state.AddResult(result)
	}

	return nil
}

func (i *Interpreter) OnCallStart(call *Call) (StepOption, error) {
	state := i.state.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}

	i.state.Push(NewCallState(call.Signature, i.state.Top().Paths(), i.cache))

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

	i.state.Push(NewCompositeState(expr.Operator, len(expr.Operands), i.state.Top().Paths()))
	return StepInto, nil
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

func (i *Interpreter) loadSet(rel *Relation) (ObjSet, error) {
	if result, ok := i.cache.LookupSet(rel); ok {
		if result == nil {
			result = NewSet[model.ObjectID]()
		}
		return result, nil
	}

	var rid dsc.RelationIdentifier
	rel.Identifier(&rid)

	relsPtr := i.pool.GetSlice()
	if err := i.loader(&rid, relsPtr); err != nil {
		return nil, err
	}

	resultSet := SetFromSlice(*relsPtr, func(rid *dsc.RelationIdentifier) model.ObjectID {
		return model.ObjectID(rid.SubjectId)
	})

	i.pool.PutSlice(relsPtr)

	i.cache.StoreSet(rel, resultSet)

	return resultSet, nil
}
