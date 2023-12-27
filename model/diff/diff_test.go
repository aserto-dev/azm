package diff_test

import (
	"strings"
	"testing"

	"github.com/aserto-dev/azm/model/diff"
	stts "github.com/aserto-dev/azm/stats"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

type (
	Stats            = stts.Stats
	ObjectTypes      = stts.ObjectTypes
	Relations        = stts.Relations
	SubjectTypes     = stts.SubjectTypes
	SubjectRelations = stts.SubjectRelations
)

func TestCanUpdateModel(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := require.New(t)
			err := canUpdate(tc.from, tc.to, tc.stats)
			tc.verify(assert, err)
		})
	}
}

var testCases = []testCase{
	{
		"no changes", baseModel, baseModel, &Stats{}, func(assert *require.Assertions, err error) { assert.NoError(err) },
	},
	{
		"delete object with no instances", baseModel, noGroupObject, &Stats{}, func(assert *require.Assertions, err error) { assert.NoError(err) },
	},
	{
		"delete object with instances", baseModel, noGroupObject,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 0}}},
		func(assert *require.Assertions, err error) {
			assert.Error(err)
			assert.ErrorContains(err, derr.ErrObjectTypeInUse.Msg("group").Error())
		},
	},
	{
		"delete object with relations", baseModel, noGroupObject,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 0, Count: 1}}},
		func(assert *require.Assertions, err error) {
			assert.Error(err)
			assert.ErrorContains(err, derr.ErrObjectTypeInUse.Msg("group").Error())
		},
	},
	{
		"delete relation with no instances", baseModel, noMemberRelation,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 0}}},
		func(assert *require.Assertions, err error) { assert.NoError(err) },
	},
	{
		"delete relation with instances", baseModel, noMemberRelation,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{"member": {Count: 1}}}}},
		func(assert *require.Assertions, err error) {
			assert.Error(err)
			assert.ErrorContains(err, derr.ErrRelationTypeInUse.Msg("group").Error())
		},
	},
	{
		"delete direct assignment with no instances", baseModel, noDirectAssignment,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{
			"member": {Count: 1, SubjectTypes: SubjectTypes{"user:*": {Count: 1}}},
		}}}},
		func(assert *require.Assertions, err error) { assert.NoError(err) },
	},
	{
		"delete direct assignment with instances", baseModel, noDirectAssignment,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{
			"member": {Count: 1, SubjectTypes: SubjectTypes{"user": {Count: 1, SubjectRelations: SubjectRelations{"": {Count: 1}}}}},
		}}}},
		func(assert *require.Assertions, err error) {
			assert.Error(err)
			assert.ErrorContains(err, derr.ErrRelationTypeInUse.Msg("group").Error())
		},
	},
	{
		"delete wildcard assignment with no instances", baseModel, noWildcardAssignment,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{
			"member": {Count: 1, SubjectTypes: SubjectTypes{"user": {Count: 1}}},
		}}}},
		func(assert *require.Assertions, err error) { assert.NoError(err) },
	},
	{
		"delete wildcard assignment with instances", baseModel, noWildcardAssignment,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{
			"member": {Count: 1, SubjectTypes: SubjectTypes{"user:*": {Count: 1, SubjectRelations: SubjectRelations{"": {Count: 1}}}}},
		}}}},
		func(assert *require.Assertions, err error) {
			assert.Error(err)
			assert.ErrorContains(err, derr.ErrRelationTypeInUse.Msg("group").Error())
		},
	},
	{
		"delete subject relation assignment with no instances", baseModel, noSubjectRelationAssignment,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{
			"member": {Count: 1, SubjectTypes: SubjectTypes{"group": {Count: 1}}},
		}}}},
		func(assert *require.Assertions, err error) { assert.NoError(err) },
	},
	{
		"delete subject relation assignment with instances", baseModel, noSubjectRelationAssignment,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{
			"member": {Count: 1, SubjectTypes: SubjectTypes{"group": {Count: 1, SubjectRelations: SubjectRelations{"member": {Count: 1}}}}},
		}}}},
		func(assert *require.Assertions, err error) {
			assert.Error(err)
			assert.ErrorContains(err, derr.ErrRelationTypeInUse.Msg("group").Error())
		},
	},
	{
		"multiple errors", baseModel, noDirectAssignment,
		&Stats{ObjectTypes: ObjectTypes{"group": {ObjCount: 1, Count: 1, Relations: Relations{
			"member": {Count: 2, SubjectTypes: SubjectTypes{
				"user":  {Count: 1, SubjectRelations: SubjectRelations{"": {Count: 1}}},
				"group": {Count: 1, SubjectRelations: SubjectRelations{"": {Count: 1}}},
			}},
		}}}},
		func(assert *require.Assertions, err error) {
			assert.Error(err)

			aerr := derr.ErrInvalidArgument
			assert.ErrorAs(err, &aerr)

			merr := aerr.Unwrap().(*multierror.Error)
			assert.Len(merr.Errors, 2)
			assert.ErrorContains(merr, derr.ErrRelationTypeInUse.Msg("group#member@user").Error())
			assert.ErrorContains(merr, derr.ErrRelationTypeInUse.Msg("group#member@group").Error())
		},
	},
}

type testCase struct {
	name   string
	from   string
	to     string
	stats  *Stats
	verify func(assert *require.Assertions, err error)
}

const (
	baseModel = `
model:
  version: 3

types:
  user:
    relations:
      manager: user

  group:
    relations:
      member: user | user:* | group#member | group
`

	noGroupObject = `
model:
  version: 3

types:
  user:
    relations:
      manager: user
`

	noMemberRelation = `
model:
  version: 3

types:
  user:
    relations:
      manager: user

  group: {}
`

	noDirectAssignment = `
model:
  version: 3

types:
  user:
    relations:
      manager: user

  group:
    relations:
      member: user:* | group#member
`
	noWildcardAssignment = `
model:
  version: 3

types:
  user:
    relations:
      manager: user

  group:
    relations:
      member: user | group#member
`

	noSubjectRelationAssignment = `
model:
  version: 3

types:
  user:
    relations:
      manager: user

  group:
    relations:
      member: user | user:* | group
`
)

func canUpdate(from, to string, stats *Stats) error {
	mFrom, err := v3.Load(strings.NewReader(from))
	if err != nil {
		return err
	}

	mTo, err := v3.Load(strings.NewReader(to))
	if err != nil {
		return err
	}

	return diff.CanUpdateModel(mFrom, mTo, stats)
}
