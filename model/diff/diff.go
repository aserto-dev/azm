package diff

import (
	"context"

	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/pkg/errors"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mock_directory_validator.go  -package=diff github.com/aserto-dev/azm/model/diff DirectoryValidator

type Diff struct {
	Added   Changes
	Removed Changes
}

type Changes struct {
	Objects   []string
	Relations map[string][]string
}

type DirectoryValidator interface {
	HasObjectInstances(ctx context.Context, objectType string) (bool, error)
	HasRelationInstances(ctx context.Context, objectType, relationName string) (bool, error)
}

func (d *Diff) Validate(ctx context.Context, dv DirectoryValidator) error {
	if err := d.validateObjectTypes(ctx, dv); err != nil {
		return err
	}

	if err := d.validateRelationsTypes(ctx, dv); err != nil {
		return err
	}

	return nil
}

func (d *Diff) validateObjectTypes(ctx context.Context, dv DirectoryValidator) error {
	for _, objType := range d.Removed.Objects {
		hasInstance, err := dv.HasObjectInstances(ctx, objType)
		if err != nil {
			return err
		}
		if hasInstance {
			return errors.Wrapf(derr.ErrObjectTypeInUse, "object type: %s", objType)
		}
	}
	return nil
}

func (d *Diff) validateRelationsTypes(ctx context.Context, dv DirectoryValidator) error {
	for objType, rels := range d.Removed.Relations {
		for _, rel := range rels {
			hasInstance, err := dv.HasRelationInstances(ctx, objType, rel)
			if err != nil {
				return err
			}
			if hasInstance {
				return errors.Wrapf(derr.ErrRelationTypeInUse, "relation type: %s; object type: %s", rel, objType)
			}
		}
	}
	return nil
}
