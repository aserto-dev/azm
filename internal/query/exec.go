package query

import (
	"slices"

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
)

type RelationLoader func(*dsc.RelationIdentifier, *Relations) error

func newRelationLoader(reader RelationReader, pool RelationPool) RelationLoader {
	return func(rid *dsc.RelationIdentifier, outRels *Relations) error {
		return reader(rid, pool, outRels)
	}
}

type Interpreter struct {
	plan       *Plan
	loader     RelationLoader
	relPool    *mempool.RelationsPool
	setPool    *ObjSetPool
	ctxFactory *ContextFactory
	cache      Cache
	context    *ds.Stack[ExecutionContext]
}

func NewInterpreter(plan *Plan, getRels RelationReader, pool *mempool.RelationsPool, setPool *ObjSetPool) *Interpreter {
	cache := Cache{}
	return &Interpreter{
		plan:       plan,
		loader:     newRelationLoader(getRels, pool),
		relPool:    pool,
		setPool:    setPool,
		ctxFactory: &ContextFactory{cache, setPool},
		cache:      cache,
	}
}

var stateSlicePool = mempool.NewSlicePool[ExecutionContext](32)

func (i *Interpreter) Run(req *dsr.CheckRequest) (ObjSet, error) {
	slicePtr := stateSlicePool.Get()
	*slicePtr = append(*slicePtr, i.ctxFactory.NewCompositeContext(Union, 1, []Scope{{ObjectID(req.ObjectId), ObjectID(req.SubjectId)}}))

	i.context = ds.AttachStack(slicePtr)
	defer func() {
		stateSlicePool.Put(i.context.Release())
	}()

	if err := i.plan.Visit(i); err != nil {
		return ObjSet{}, err
	}

	if i.context.Len() != 1 {
		return ObjSet{}, errors.Wrap(ErrInterpreter, "unbalanced stack")
	}

	i.cache.Clear(i.setPool)

	return i.context.Top().Result(), nil
}

func (i *Interpreter) OnLoad(expr *Load) error {
	state := i.context.Top()

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

		result, err := i.loadSet(&Relation{RelationType: expr.RelationType, Scope: &scope})
		if err != nil {
			return err
		}

		state.AddSet(result)
	}

	return nil
}

func (i *Interpreter) OnPipeStart(pipe *Pipe) (StepOption, error) {
	state := i.context.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}
	i.context.Push(i.ctxFactory.NewPipeContext(state.Scopes()))
	return StepInto, nil
}

func (i *Interpreter) OnPipeEnd(_ *Pipe) {
	i.rollupResult()
}

func (i *Interpreter) OnCallStart(call *Call) (StepOption, error) {
	state := i.context.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}

	i.context.Push(i.ctxFactory.NewCallContext(call.Signature, i.context.Top().Scopes()))

	return StepInto, nil
}

func (i *Interpreter) OnCallEnd(_ *Call) {
	i.rollupResult()
}

func (i *Interpreter) OnCompositeStart(expr *Composite) (StepOption, error) {
	state := i.context.Top()
	if state.ShortCircuit() {
		return StepOver, nil
	}

	i.context.Push(i.ctxFactory.NewCompositeContext(expr.Operator, len(expr.Operands), i.context.Top().Scopes()))
	return StepInto, nil
}

func (i *Interpreter) OnCompositeEnd(_ *Composite) {
	i.rollupResult()
}

func (i *Interpreter) rollupResult() {
	if i.context.Len() > 1 {
		state := i.context.Pop()
		result := state.Result()
		i.context.Top().AddSet(result)
		i.setPool.PutSet(result)
	}
}

func (i *Interpreter) loadSet(rel *Relation) (ObjSet, error) {
	if result, ok := i.cache.LookupSet(rel); ok {
		if result == nil {
			result = &ObjSet{}
		}
		return *result, nil
	}

	var rid dsc.RelationIdentifier
	rel.Identifier(&rid)

	relsPtr := i.relPool.GetSlice()
	if err := i.loader(&rid, relsPtr); err != nil {
		return ObjSet{}, err
	}

	resultSet := i.setPool.GetSet()
	resultSet.Add(ds.TransformIter(slices.Values(*relsPtr), func(rid *dsc.RelationIdentifier) model.ObjectID {
		return model.ObjectID(rid.SubjectId)
	}))

	i.relPool.PutSlice(relsPtr)

	i.cache.StoreSet(rel, resultSet)

	return resultSet, nil
}
