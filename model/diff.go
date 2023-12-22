package model

import (
	stats "github.com/aserto-dev/azm/stats"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock_instances.go  -package=model github.com/aserto-dev/azm/model Instances

type Diff struct {
	Added   Changes
	Removed Changes
}

type Changes struct {
	Objects   []ObjectName
	Relations map[ObjectName]map[RelationName][]string
}
type Instances interface {
	GetStats() (*stats.Stats, error)
}

func (d *Diff) Validate(dv Instances) error {
	var errs error

	if len(d.Removed.Objects) == 0 && len(d.Removed.Relations) == 0 {
		return errs
	}

	sts, err := dv.GetStats()
	if err != nil {
		return err
	}

	if len(d.Removed.Objects) > 0 {
		err := d.validateObjectTypes(sts)
		errs = multierror.Append(errs, err)
	}

	if len(d.Removed.Relations) > 0 {
		err := d.validateRelationsTypes(sts)
		errs = multierror.Append(errs, err)
	}

	if merr, ok := errs.(*multierror.Error); ok && len(merr.Errors) > 0 {
		return errs
	}

	return nil
}

func (d *Diff) validateObjectTypes(sts *stats.Stats) error {
	var errs error
	for _, objType := range d.Removed.Objects {
		if sts.ObjectTypes[objType].ObjCount > 0 {
			errs = multierror.Append(errs, errors.Wrapf(derr.ErrObjectTypeInUse, "object type [%s]", objType))
		}

		for relName, relStats := range sts.ObjectTypes[objType].Relations {
			if relStats.Count > 0 {
				errs = multierror.Append(errs, errors.Wrapf(derr.ErrRelationTypeInUse, "object type [%s], relation type [%s]", objType, relName))
			}
		}

	}
	return errs
}

func (d *Diff) validateRelationsTypes(sts *stats.Stats) error {
	var errs error
	for objType, rels := range d.Removed.Relations {
		for rel, subjs := range rels {
			if len(subjs) == 0 {
				if sts.ObjectTypes[objType].Relations[rel].Count > 0 {
					errs = multierror.Append(errs, errors.Wrapf(derr.ErrRelationTypeInUse, "object type [%s], relation type [%s]", objType, rel))
				}
			} else {
				for _, sub := range subjs {
					if sts.ObjectTypes[objType].Relations[rel].SubjectTypes[ObjectName(sub)].Count > 0 {
						errs = multierror.Append(errs, errors.Wrapf(derr.ErrRelationTypeInUse, "object type [%s], relation type [%s], subject type [%s]", objType, rel, sub))
					}
				}
			}

		}
	}
	return errs
}
