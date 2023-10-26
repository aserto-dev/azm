package diff

import (
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock_instances.go  -package=diff github.com/aserto-dev/azm/model/diff Instances

type Diff struct {
	Added   Changes
	Removed Changes
}

type Changes struct {
	Objects   []string
	Relations map[string][]string
}

// Only the types of the relation instances are needed.
type RelationKind struct {
	Object          string
	Relation        string
	Subject         string
	SubjectRelation string
}

type Instances interface {
	ObjectTypes() ([]string, error)
	RelationTypes() ([]*RelationKind, error)
}

func (d *Diff) Validate(dv Instances) error {
	var errs error
	var rels []*RelationKind
	if len(d.Removed.Objects) > 0 {
		objs, err := dv.ObjectTypes()
		if err != nil {
			return err
		}

		rels, err = dv.RelationTypes()
		if err != nil {
			return err
		}

		err = d.validateObjectTypes(objs, rels)
		errs = multierror.Append(errs, err)
	}

	if len(d.Removed.Relations) > 0 {
		var err error
		if len(rels) == 0 {
			rels, err = dv.RelationTypes()
			if err != nil {
				return err
			}

		}
		err = d.validateRelationsTypes(rels)
		errs = multierror.Append(errs, err)
	}

	if merr, ok := errs.(*multierror.Error); ok && len(merr.Errors) > 0 {
		return errs
	}

	return nil
}

func (d *Diff) validateObjectTypes(objs []string, rels []*RelationKind) error {
	var errs error
	for _, objType := range d.Removed.Objects {
		_, found := lo.Find(objs, func(obj string) bool { return obj == objType })
		if found {
			errs = multierror.Append(errs, errors.Wrapf(derr.ErrObjectTypeInUse, "object type [%s]", objType))
		}
		rel, found := lo.Find(rels, func(rel *RelationKind) bool { return rel.Object == objType || rel.Subject == objType })
		if found {
			errs = multierror.Append(errs, errors.Wrapf(derr.ErrRelationTypeInUse, "object type [%s], relation type [%s]", objType, rel.Relation))
		}
	}
	return errs
}

func (d *Diff) validateRelationsTypes(relations []*RelationKind) error {
	var errs error
	for objType, rels := range d.Removed.Relations {
		for _, rel := range rels {
			_, found := lo.Find(relations, func(rl *RelationKind) bool {
				return (rl.Object == objType && rl.Relation == rel) || (rl.Subject == objType && rl.SubjectRelation == rel)
			})
			if found {
				errs = multierror.Append(errs, errors.Wrapf(derr.ErrRelationTypeInUse, "object type [%s], relation type [%s]", objType, rel))
			}
		}
	}
	return errs
}
