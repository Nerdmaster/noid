package noid

import (
	"strings"
	"errors"
)

type Ordering int

const (
	Random Ordering = iota
	SequentialLimited
	SequentialUnlimited
)

type Template struct {
	prefix string
	ordering Ordering
	mask string
	hasCheckDigit bool
}

func NewTemplate(template string) (*Template, error) {
	var suffix string
	var err error

	// You know what's hip and cool these days?  Storing values immediately on
	// instantiation when said values are essentially static, read-only data
	t := &Template{}
	t.prefix, suffix = splitTemplateString(template)
	t.hasCheckDigit, suffix = getCheckDigitFromSuffix(suffix)
	t.ordering, err = getOrderingFromChar(suffix[0])

	if err != nil {
		return nil, err
	}

	t.mask = suffix[1:]

	return t, nil
}

// Returns the prefix and suffix parts of a template string, defaulting to a
// prefix of "" when no period is in the string
func splitTemplateString(s string) (string, string) {
	parts := strings.Split(s, ".")

	if len(parts) == 1 {
		return "", s
	}

	return parts[0], parts[1]
}

// Returns whether or not the final character is a check digit ("k") as well as
// the new suffix (in the case of no check digit, the suffix returned is the
// same as was passed in)
func getCheckDigitFromSuffix(suffix string) (bool, string) {
	last := len(suffix) - 1
	if suffix[last] == 'k' {
		return true, suffix[0:last]
	}

	return false, suffix
}

func getOrderingFromChar(c byte) (Ordering, error) {
	var err error
	var order Ordering

	err = nil

	switch c {
		case 'r': order = Random
		case 's': order = SequentialLimited
		case 'z': order = SequentialUnlimited
		default: err = errors.New("Ordering character must be 'r', 's', or 'z'")
	}

	return order, err
}
