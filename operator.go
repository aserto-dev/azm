package azm

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Operator int

const (
	None      Operator = iota // NONE
	This                      // THIS
	Union                     // OR
	Intersect                 // AND
	Exclude                   // NOT (slot 0 = base, slot 1 = subtraction)
	Function                  // External functional call ($rego:<fq_rule> or $wasm:<fq_fn>)
)

const (
	none      string = "none"
	this      string = "this"
	union     string = "union"
	intersect string = "intersect"
	exclude   string = "exclude"
	function  string = "function"
)

var (
	operationToString = map[Operator]string{
		None:      none,
		This:      this,
		Union:     union,
		Intersect: intersect,
		Exclude:   exclude,
		Function:  function,
	}
	operationToEnum = map[string]Operator{
		none:      None,
		this:      This,
		union:     Union,
		intersect: Intersect,
		exclude:   Exclude,
		function:  Function,
	}
)

// String, lookup string representation of Operator.
func (o Operator) String() string {
	s, ok := operationToString[o]
	if ok {
		return s
	}
	return ""
}

// MarshalJSON Operator enum to string representation.
func (o *Operator) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

// UnmarshalJSON Operator string representation to enum instance.
func (o *Operator) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	e, ok := operationToEnum[s]
	if !ok {
		return errors.Wrapf(ErrOperationNotFound, "%s", s)
	}

	*o = e

	return nil
}

// MarshalYAML Operator enum to string representation.
func (o Operator) MarshalYAML() (interface{}, error) {
	s := o.String()
	return s, nil
}

// Unmarshal Operator string representation to enum instance.
func (o *Operator) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		*o = None
	}

	e, ok := operationToEnum[s]
	if !ok {
		return errors.Wrapf(ErrOperationNotFound, "%s", s)
	}

	*o = e

	return nil
}
