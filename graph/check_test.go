package graph_test

import (
	"strings"
	"testing"

	azmgraph "github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/internal/mempool"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		check    string
		expected bool
	}{
		// Relations
		{"doc:doc1#owner@user:user1", false},
		{"doc:doc1#viewer@user:user1", true},
		{"doc:doc2#viewer@user:user1", true},
		{"doc:doc2#viewer@user:userX", true},
		{"doc:doc1#viewer@user:user2", true},
		{"doc:doc1#viewer@user:user3", true},
		{"doc:doc1#viewer@group:d1_viewers", false},

		{"group:yin#member@user:yin_user", true},
		{"group:yin#member@user:yang_user", true},
		{"group:yang#member@user:yin_user", true},
		{"group:yang#member@user:yang_user", true},

		{"group:alpha#member@user:user1", false},

		// Permissions
		{"doc:doc1#can_change_owner@user:d1_owner", true},
		{"doc:doc1#can_change_owner@user:user1", false},
		{"doc:doc1#can_change_owner@user:userX", false},

		{"doc:doc1#can_read@user:d1_owner", true},
		{"doc:doc1#can_read@user:f1_owner", true},
		{"doc:doc1#can_read@user:user1", true},
		{"doc:doc1#can_read@user:f1_viewer", true},
		{"doc:doc1#can_read@user:userX", false},

		{"doc:doc1#can_write@user:d1_owner", true},
		{"doc:doc1#can_write@user:f1_owner", true},
		{"doc:doc1#can_write@user:user2", false},

		{"folder:folder1#owner@user:f1_owner", true},
		{"folder:folder1#can_create_file@user:f1_owner", true},
		{"folder:folder1#can_share@user:f1_owner", true},

		// intersection
		{"doc:doc1#can_share@user:d1_owner", false},
		{"doc:doc1#can_share@user:f1_owner", true},

		// negation
		{"folder:folder1#can_read@user:f1_owner", true},
		{"doc:doc1#viewer@user:f1_owner", false},
		{"doc:doc1#can_invite@user:f1_owner", true},

		// cycles
		{"cycle:loop#can_delete@user:loop_owner", true},
		{"cycle:loop#can_delete@user:user1", false},
	}

	m, err := v3.LoadFile("./check_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, m)

	pool := mempool.NewSlicePool[*dsc.Relation]()

	for _, test := range tests {
		t.Run(test.check, func(tt *testing.T) {
			assert := assert.New(tt)

			checker := azmgraph.NewCheck(m, checkReq(test.check), rels.GetRelations, pool)

			res, err := checker.Check()
			assert.NoError(err)
			tt.Log("trace:\n", strings.Join(checker.Trace(), "\n"))
			assert.Equal(test.expected, res)
		})
	}

}

var rels = NewRelationsReader(
	"folder:folder1#owner@user:f1_owner",
	"folder:folder1#viewer@group:f1_viewers#member",
	"group:f1_viewers#member@user:f1_viewer",
	"doc:doc1#parent@folder:folder1",
	"doc:doc1#owner@user:d1_owner",
	"doc:doc1#viewer@group:d1_viewers#member",
	"doc:doc1#viewer@user:user1",
	"group:d1_viewers#member@user:user2",
	"doc:doc2#viewer@user:*",

	"group:d1_viewers#member@group:d1_subviewers#member",
	"group:d1_subviewers#member@user:user3",

	// mutually recursive groups with users
	"group:yin#member@group:yang#member",
	"group:yang#member@group:yin#member",
	"group:yin#member@user:yin_user",
	"group:yang#member@user:yang_user",

	// mutually recursive groups with no users
	"group:alpha#member@group:omega#member",
	"group:omega#member@group:alpha#member",

	// cyclical permissions
	"cycle:loop#parent@cycle:loop",
	"cycle:loop#owner@user:loop_owner",
)
