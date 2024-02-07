package check

import (
	"fmt"
	"strings"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/samber/lo"
)

type checkStatus int

func (s checkStatus) String() string {
	switch s {
	case checkStatusUnknown:
		return "unknown"
	case checkStatusPending:
		return "?"
	case checkStatusTrue:
		return "true"
	case checkStatusFalse:
		return "false"
	default:
		return fmt.Sprintf("invalid: %d", s)
	}
}

const (
	checkStatusUnknown checkStatus = iota
	checkStatusPending
	checkStatusTrue
	checkStatusFalse
)

type CheckResults []CheckParams

func (r CheckResults) status() checkStatus {
	switch {
	case r == nil:
		return checkStatusPending
	case len(r) == 0:
		return checkStatusFalse
	default:
		return checkStatusTrue
	}
}

func (r CheckResults) addResult(rels ...*dsc.Relation) CheckResults {
	results := lo.Map(rels, func(rel *dsc.Relation, _ int) CheckParams {
		return CheckParams{
			OT:   model.ObjectName(rel.ObjectType),
			OID:  ObjectID(rel.ObjectId),
			Rel:  model.RelationName(rel.Relation),
			ST:   model.ObjectName(rel.SubjectType),
			SID:  ObjectID(rel.SubjectId),
			SRel: model.RelationName(rel.SubjectRelation),
		}
	})
	return r.append(results)
}

func (r CheckResults) append(results CheckResults) CheckResults {
	return lo.Uniq(append(r, results...))
}

type checkCall struct {
	*CheckParams
	results CheckResults
}

// checkMemo tracks pending checks to detect cycles, and caches the results of completed checks.
type checkMemo struct {
	memo    map[CheckParams]CheckResults
	history []*checkCall
}

func newCheckMemo(trace bool) *checkMemo {
	return &checkMemo{
		memo:    map[CheckParams]CheckResults{},
		history: lo.Ternary(trace, []*checkCall{}, nil),
	}
}

// MarkVisited returns the status of a check. If the check has not been visited, it is marked as pending.
func (m *checkMemo) MarkVisited(params *CheckParams) checkStatus {
	prior, ok := m.memo[*params]
	current := prior
	if !ok {
		var current CheckResults
		m.memo[*params] = current
	}

	m.trace(params, current)

	switch {
	case !ok:
		return checkStatusUnknown
	case prior == nil:
		return checkStatusPending
	case len(prior) == 0:
		return checkStatusFalse
	default:
		return checkStatusTrue
	}
}

// MarkComplete records the result of a check.
func (m *checkMemo) MarkComplete(params *CheckParams, results CheckResults) {
	m.memo[*params] = results

	m.trace(params, results)
}

func (m *checkMemo) Results(params *CheckParams) CheckResults {
	return m.memo[*params]
}

func (m *checkMemo) Trace() []string {
	if m.history == nil {
		return []string{}
	}

	callstack := []string{}

	return lo.Map(m.history, func(c *checkCall, _ int) string {
		call := c.String()
		status := c.results.status()

		if len(callstack) > 0 && callstack[len(callstack)-1] == call && status != checkStatusPending {
			callstack = callstack[:len(callstack)-1]
		}

		s := fmt.Sprintf("%s%s = %s", strings.Repeat("  ", len(callstack)), call, status)

		if status == checkStatusPending {
			callstack = append(callstack, call)
		}

		return s
	})
}

func (m *checkMemo) trace(params *CheckParams, results CheckResults) {
	if m.history != nil {
		m.history = append(m.history, &checkCall{params, results})
	}
}
