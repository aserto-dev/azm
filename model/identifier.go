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
// - have a minimum length of 3 characters
// - have a maximum length of 64 characters
// - start with a character (a-zA-Z)
// - end with a character of a digit (a-zA-Z0-9)
// - can contain dots, underscores and dashes, between the first and last position.
type Identifier string

var ErrInvalidIdentifier = errors.New("invalid identifier")

var (
	reIdentifier         = regexp.MustCompile(`(?m)^[a-zA-Z][a-zA-Z0-9._-]{1,62}[a-zA-Z0-9]$`)
	msgInvalidIdentifier = "must start with a letter, can contain mixed case letters, digits, dots, underscores, and dashes, and must end with a letter or digit"
)

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
