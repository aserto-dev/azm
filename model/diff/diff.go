package diff

import (
	"context"

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
	ObjectsExist(ctx context.Context, objectType string) (bool, error)
	RelationsExist(ctx context.Context, objectType, relationName string) (bool, error)
}

func (d *Diff) Validate(ctx context.Context, dv Instances) error {
	var errs error
	if err := d.validateObjectTypes(ctx, dv); err != nil {
		errs = multierror.Append(errs, err)
	}

	if err := d.validateRelationsTypes(ctx, dv); err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs
}

func (d *Diff) validateObjectTypes(ctx context.Context, dv Instances) error {
	var errs error
	for _, objType := range d.Removed.Objects {
		hasInstance, err := dv.ObjectsExist(ctx, objType)
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

func (d *Diff) validateRelationsTypes(ctx context.Context, dv Instances) error {
	var errs error
	for objType, rels := range d.Removed.Relations {
		for _, rel := range rels {
			hasInstance, err := dv.RelationsExist(ctx, objType, rel)
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
