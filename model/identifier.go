package model

import (
	"errors"
	"regexp"
	"strings"
)

// Identifier is the string representation of an object, relation and permission type name.
//
// Identifiers are bounded by the underlying defined regex definition (reIdentifier).
//
// An identifier MUST be:
// - all lowercase characters
// - have a minimum length of 3 characters
// - have a maximum length of 64 characters
// - start with a character (a-z)
// - end with a character of a digit (a-z0-9)
// - can contain dots, underscores and dashes, between the first and last position.
type Identifier string

var ErrInvalidIdentifier = errors.New("invalid identifier")

var reIdentifier = regexp.MustCompile(`(?m)^[a-z][a-z0-9._-]{1,62}[a-z0-9]$`)
var msgInvalidIdentifier = "must be all lowercase, start with a letter, can contain letters, digits, dots, underscores, and dashes, and must end with a letter or digit"

func (i Identifier) Valid() bool {
	return reIdentifier.MatchString(string(i))
}

func IsValidIdentifier(in string) bool {
	return reIdentifier.MatchString(in)
}

func NormalizeIdentifier(in string) (string, error) {
	if IsValidIdentifier(in) {
		return in, nil
	}

	// always lowercase the identifier
	in = strings.ToLower(in)
	if IsValidIdentifier(in) {
		return in, nil
	}

	return in, ErrInvalidIdentifier
}
