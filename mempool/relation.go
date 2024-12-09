package mempool

import (
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RelationsPool = CollectionPool[*dsc.Relation]

func NewRelationsPool() *RelationsPool {
	return NewCollectionPool[*dsc.Relation](NewRelationAllocator())
}

type RelationAllocator struct {
	tsPool *Pool[*timestamppb.Timestamp]
}

func NewRelationAllocator() *RelationAllocator {
	return &RelationAllocator{
		tsPool: NewPool[*timestamppb.Timestamp](
			func() *timestamppb.Timestamp {
				return new(timestamppb.Timestamp)
			}),
	}
}

func (ra *RelationAllocator) New() *dsc.Relation {
	rel := dsc.RelationFromVTPool()
	rel.CreatedAt = ra.tsPool.Get()
	rel.UpdatedAt = ra.tsPool.Get()
	return rel
}

func (ra *RelationAllocator) Reset(rel *dsc.Relation) {
	if rel.CreatedAt != nil {
		rel.CreatedAt.Reset()
		ra.tsPool.Put(rel.CreatedAt)
	}
	if rel.UpdatedAt != nil {
		rel.UpdatedAt.Reset()
		ra.tsPool.Put(rel.UpdatedAt)
	}

	rel.ReturnToVTPool()
}
