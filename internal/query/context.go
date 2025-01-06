package query

import (
	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/model"
)

type (
	ObjSet     = ds.Set[model.ObjectID]
	ObjSetPool = ds.SetPool[model.ObjectID]
)

type Scope struct {
	OID model.ObjectID
	SID model.ObjectID
}

type ExecutionContext interface {
	AddSet(ObjSet)
	ShortCircuit() bool
	Scopes() []Scope
	Result() ObjSet
}

type ContextFactory struct {
	cache   Cache
	setPool *ObjSetPool
}

func (f *ContextFactory) NewCallContext(sig *RelationType, scopes []Scope) *CallContext {
	return newCallContext(sig, scopes, f.setPool.GetSet(), f.cache)
}

func (f *ContextFactory) NewCompositeContext(op Operator, size int, scopes []Scope) *CompositeContext {
	return newCompositeContext(op, size, scopes, f.setPool.GetSet())
}

func (f *ContextFactory) NewPipeContext(scopes []Scope) *PipeContext {
	return newPipeContext(scopes, f.setPool.GetSet())

}
