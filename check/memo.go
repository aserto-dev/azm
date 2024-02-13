package check

import (
	"fmt"
	"strings"

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

type checkCall struct {
	*checkParams
	status checkStatus
}

// checkMemo tracks pending checks to detect cycles, and caches the results of completed checks.
type checkMemo struct {
	memo    map[checkParams]checkStatus
	history []*checkCall
}

func newCheckMemo(trace bool) *checkMemo {
	return &checkMemo{
		memo:    map[checkParams]checkStatus{},
		history: lo.Ternary(trace, []*checkCall{}, nil),
	}
}

// MarkVisited returns the status of a check. If the check has not been visited, it is marked as pending.
func (m *checkMemo) MarkVisited(params *checkParams) checkStatus {
	prior := m.memo[*params]
	current := prior
	if prior == checkStatusUnknown {
		current = checkStatusPending
		m.memo[*params] = current
	}

	m.trace(params, current)

	return prior
}

// MarkComplete records the result of a check.
func (m *checkMemo) MarkComplete(params *checkParams, checkResult bool) {
	status := checkStatusFalse
	if checkResult {
		status = checkStatusTrue
	}
	m.memo[*params] = status

	m.trace(params, status)
}

func (m *checkMemo) Trace() []string {
	if m.history == nil {
		return []string{}
	}

	callstack := []string{}

	return lo.Map(m.history, func(c *checkCall, _ int) string {
		call := c.String()
		result := c.status.String()

		if len(callstack) > 0 && callstack[len(callstack)-1] == call && c.status != checkStatusPending {
			callstack = callstack[:len(callstack)-1]
		}

		s := fmt.Sprintf("%s%s = %s", strings.Repeat("  ", len(callstack)), call, result)

		if c.status == checkStatusPending {
			callstack = append(callstack, call)
		}

		return s
	})
}

func (m *checkMemo) trace(params *checkParams, status checkStatus) {
	if m.history != nil {
		m.history = append(m.history, &checkCall{params, status})
	}
}
