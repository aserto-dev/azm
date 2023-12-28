package diff

import (
	"github.com/aserto-dev/azm/model"
	stts "github.com/aserto-dev/azm/stats"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
)

type (
	Stats        = stts.Stats
	ObjectName   = model.ObjectName
	RelationName = model.RelationName
)

func CanUpdateModel(cur, next *model.Model, stats *Stats) error {
	d := calculateDelta(cur, next)
	if len(d) == 0 {
		return nil
	}

	var errs error
	for on, rd := range d {
		if len(rd) == 0 && stats.ObjectRefCount(on) > 0 {
			// The object has been removed but there are still instances or relations.
			errs = multierror.Append(errs, derr.ErrObjectTypeInUse.Msg(on.String()))
			continue
		}

		for rn, rel := range rd {
			if len(rel) == 0 && stats.RelationRefCount(on, rn) > 0 {
				// The relation has been removed but there are still instances.
				errs = multierror.Append(errs, derr.ErrRelationTypeInUse.Msgf("%s#%s", on, rn))
				continue
			}

			for ref := range rel {
				sn, sr := ref.Object, ref.Relation
				if ref.IsWildcard() {
					sn += ":*"
					sr = ""
				}
				if stats.RelationSubjectCount(on, rn, sn, sr) > 0 {
					// The relation hasn't been removed, but some of its subjects have.
					errs = multierror.Append(errs, derr.ErrRelationTypeInUse.Msgf("%s#%s@%s", on, rn, &ref))
				}
			}
		}
	}

	if errs != nil {
		return derr.ErrInvalidArgument.Err(errs)
	}

	return nil
}

func calculateDelta(from, sub *model.Model) delta {
	chgs := delta{}

	if from == nil {
		return chgs
	}

	if sub == nil {
		return lo.MapValues(from.Objects, func(_ *model.Object, _ model.ObjectName) dRelations { return dRelations{} })
	}

	for objName, obj := range from.Objects {
		subObj := sub.Objects[objName]
		if subObj == nil {
			chgs[objName] = dRelations{}
			continue
		}

		relsDiff := dRelations{}

		for relname, rel := range obj.Relations {
			if subObj.Relations[relname] == nil {
				relsDiff[relname] = dRelation{}
				continue
			}

			relDiff, _ := DifferencePtr(rel.Union, sub.Objects[objName].Relations[relname].Union)
			relsDiff[relname] = lo.Associate(relDiff, func(r *model.RelationRef) (model.RelationRef, struct{}) { return *r, struct{}{} })
		}

		if len(relsDiff) > 0 {
			chgs[objName] = relsDiff
		}
	}

	return chgs
}

type (
	delta      map[ObjectName]dRelations
	dRelations map[RelationName]dRelation
	dRelation  map[model.RelationRef]struct{}
)

// Similar to lo.Difference but for slices of pointers.
func DifferencePtr[T comparable](list1, list2 []*T) ([]*T, []*T) {
	left := []*T{}
	right := []*T{}

	seenLeft := map[T]struct{}{}
	seenRight := map[T]struct{}{}

	for _, elem := range list1 {
		seenLeft[*elem] = struct{}{}
	}

	for _, elem := range list2 {
		seenRight[*elem] = struct{}{}
	}

	for _, elem := range list1 {
		if _, ok := seenRight[*elem]; !ok {
			left = append(left, elem)
		}
	}

	for _, elem := range list2 {
		if _, ok := seenLeft[*elem]; !ok {
			right = append(right, elem)
		}
	}

	return left, right
}
