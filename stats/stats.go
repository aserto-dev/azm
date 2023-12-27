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
