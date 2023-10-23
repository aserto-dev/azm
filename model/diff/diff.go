package diff

import (
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
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

type Instances interface {
	ObjectsExist(objectType string) (bool, error)
	RelationsExist(objectType, relationName string) (bool, error)
}

func (d *Diff) Validate(dv Instances) error {
	var errs error
	if err := d.validateObjectTypes(dv); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := d.validateRelationsTypes(dv); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs
}

func (d *Diff) validateObjectTypes(dv Instances) error {
	var errs error
	for _, objType := range d.Removed.Objects {
		hasInstance, err := dv.ObjectsExist(objType)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		if hasInstance {
			errs = multierror.Append(errs, errors.Wrapf(derr.ErrObjectTypeInUse, "object type: %s", objType))
		}
	}
	return errs
}

func (d *Diff) validateRelationsTypes(dv Instances) error {
	var errs error
	for objType, rels := range d.Removed.Relations {
		for _, rel := range rels {
			hasInstance, err := dv.RelationsExist(objType, rel)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			if hasInstance {
				errs = multierror.Append(errs, errors.Wrapf(derr.ErrRelationTypeInUse, "relation type: %s; object type: %s", rel, objType))
			}
		}
	}
	return errs
}
