package stats

import (
	"sync/atomic"

	"github.com/aserto-dev/azm/types"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
)

type Stats struct {
	ObjectTypes ObjectTypes `json:"object_types,omitempty"`
}

type ObjectTypes map[types.ObjectName]struct {
	ObjCount  int32     `json:"_obj_count,omitempty"`
	Count     int32     `json:"_count,omitempty"`
	Relations Relations `json:"relations,omitempty"`
}

type Relations map[types.RelationName]struct {
	Count        int32        `json:"_count,omitempty"`
	SubjectTypes SubjectTypes `json:"subject_types,omitempty"`
}

type SubjectTypes map[types.ObjectName]struct {
	Count            int32            `json:"_count,omitempty"`
	SubjectRelations SubjectRelations `json:"subject_relations,omitempty"`
}

type SubjectRelations map[types.RelationName]struct {
	Count int32 `json:"_count,omitempty"`
}

func (s *Stats) CountObject(obj *dsc.Object) {
	ot, ok := s.ObjectTypes[types.ObjectName(obj.Type)]
	if !ok {
		atomic.StoreInt32(&ot.ObjCount, 0)
		if ot.Relations == nil {
			ot.Relations = Relations{}
		}
	}

	atomic.AddInt32(&ot.ObjCount, 1)

	s.ObjectTypes[types.ObjectName(obj.Type)] = ot
}

func (s *Stats) CountRelation(rel *dsc.Relation) {
	objType := types.ObjectName(rel.ObjectType)
	relation := types.RelationName(rel.Relation)
	subType := types.ObjectName(rel.SubjectType)
	subRel := types.RelationName(rel.SubjectRelation)

	// object_types
	ot, ok := s.ObjectTypes[objType]
	if !ok {
		atomic.StoreInt32(&ot.Count, 0)
	}

	if ot.Relations == nil {
		ot.Relations = Relations{}
	}

	atomic.AddInt32(&ot.Count, 1)
	s.ObjectTypes[objType] = ot

	// relations
	re, ok := ot.Relations[relation]
	if !ok {
		atomic.StoreInt32(&re.Count, 0)
		re.SubjectTypes = SubjectTypes{}
	}

	atomic.AddInt32(&re.Count, 1)
	ot.Relations[relation] = re

	// subject_types
	st, ok := re.SubjectTypes[subType]
	if !ok {
		atomic.StoreInt32(&st.Count, 0)
		st.SubjectRelations = SubjectRelations{}
	}

	atomic.AddInt32(&st.Count, 1)
	re.SubjectTypes[subType] = st

	// subject_relations
	if subRel != "" {
		sr, ok := st.SubjectRelations[subRel]
		if !ok {
			atomic.StoreInt32(&sr.Count, 0)
		}

		atomic.AddInt32(&sr.Count, 1)
		st.SubjectRelations[subRel] = sr
	}
}
