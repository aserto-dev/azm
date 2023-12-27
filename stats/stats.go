package stats

import (
	"github.com/aserto-dev/azm/model"
)

type Stats struct {
	ObjectTypes ObjectTypes `json:"object_types,omitempty"`
}

type ObjectTypes map[model.ObjectName]struct {
	ObjCount  int32     `json:"_obj_count,omitempty"`
	Count     int32     `json:"_count,omitempty"`
	Relations Relations `json:"relations,omitempty"`
}

type Relations map[model.RelationName]struct {
	Count        int32        `json:"_count,omitempty"`
	SubjectTypes SubjectTypes `json:"subject_types,omitempty"`
}

type SubjectTypes map[model.ObjectName]struct {
	Count            int32            `json:"_count,omitempty"`
	SubjectRelations SubjectRelations `json:"subject_relations,omitempty"`
}

type SubjectRelations map[model.RelationName]struct {
	Count int32 `json:"_count,omitempty"`
}

func (s *Stats) ObjectRefCount(on model.ObjectName) int32 {
	ot := s.ObjectTypes[on]
	return ot.Count + ot.ObjCount
}

func (s *Stats) RelationRefCount(on model.ObjectName, rn model.RelationName) int32 {
	return s.ObjectTypes[on].Relations[rn].Count
}

func (s *Stats) RelationSubjectCount(on model.ObjectName, rn model.RelationName, sn model.ObjectName, sr model.RelationName) int32 {
	st := s.ObjectTypes[on].Relations[rn].SubjectTypes[sn]
	if sr != "" {
		return st.SubjectRelations[sr].Count
	}

	subjCount := int32(0)
	for _, subj := range st.SubjectRelations {
		subjCount += subj.Count
	}
	return st.Count - subjCount
}
