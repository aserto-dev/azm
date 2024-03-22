package graph

import (
	"fmt"
	"strings"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
)

type ObjectSearch struct {
	graphSearch
}

func NewObjectSearch(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader) *SubjectSearch {
	im := m.Invert()
	if err := im.Validate(); err != nil {
		panic(err)
	}

	return &SubjectSearch{graphSearch{
		m:       im,
		params:  invertGetGraphRequest(im, req),
		getRels: invertedRelationReader(reader),
		memo:    newSearchMemo(req.Trace),
		explain: req.Explain,
	}}
}

func invertGetGraphRequest(im *model.Model, req *dsr.GetGraphRequest) *relation {
	iReq := &relation{
		ot:  model.ObjectName(req.SubjectType),
		oid: ObjectID(req.SubjectId),
		rel: model.RelationName(fmt.Sprintf("%s_%s", req.ObjectType, req.Relation)),
		st:  model.ObjectName(req.ObjectType),
		sid: ObjectID(req.ObjectId),
		// TODO: what do we do with subject relations
		// srel: model.RelationName(req.SubjectRelation),
	}

	o := im.Objects[model.ObjectName(iReq.ot)]
	if o.HasRelation(model.RelationName(iReq.rel)) && o.HasPermission(model.RelationName("r_"+iReq.rel)) {
		iReq.rel = model.RelationName("r_" + iReq.rel)
	}

	return iReq
}

func invertedRelationReader(reader RelationReader) RelationReader {
	return func(r *dsc.Relation) ([]*dsc.Relation, error) {
		x := strings.SplitN(r.Relation, "_", 2)

		rr := &dsc.Relation{
			ObjectType:  r.SubjectType,
			ObjectId:    r.SubjectId,
			Relation:    x[1],
			SubjectType: r.ObjectType,
			SubjectId:   r.ObjectId,
		}

		res, err := reader(rr)
		if err != nil {
			return nil, err
		}

		return lo.Map(res, func(r *dsc.Relation, _ int) *dsc.Relation {
			return &dsc.Relation{
				ObjectType:  r.SubjectType,
				ObjectId:    r.SubjectId,
				Relation:    r.Relation,
				SubjectType: r.ObjectType,
				SubjectId:   r.ObjectId,
			}
		}), nil
	}
}
