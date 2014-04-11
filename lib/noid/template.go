package noid

import (
	"strings"
	"errors"
	"fmt"
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

// 14-character masks can theoretically result in a string that can't be
// represented by an int64 index, and for now I want all noids to be easily
// converted to and from a raw number.  13 characters, if all are using the
// "extended" set, is still enough to hold over 10 quintillion (10 billion
// billion) noids.
const MaxMaskLength = 13

func NewTemplate(template string) (*Template, error) {
	var suffix string

	// You know what's hip and cool these days?  Storing values immediately on
	// instantiation when said values are essentially static, read-only data
	t := &Template{}
	t.prefix, suffix = splitTemplateString(template)

	last := len(suffix) - 1
	if suffix[last] == 'k' {
		t.hasCheckDigit = true
		suffix = suffix[0:last]
	}

	switch suffix[0] {
		case 'r': t.ordering = Random
		case 's': t.ordering = SequentialLimited
		case 'z': t.ordering = SequentialUnlimited
		default: return nil, errors.New("Ordering character must be 'r', 's', or 'z'")
	}

	t.mask = suffix[1:]

	if len(t.mask) > MaxMaskLength {
		return nil, errors.New(fmt.Sprintf("Mask cannot be more than %d characters", MaxMaskLength))
	}

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
