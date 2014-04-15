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

const Digits = "0123456789"
const ExtendedDigits = "0123456789bcdfghjkmnpqrstvwxz"

type Template struct {
	prefix string
	ordering Ordering
	mask string
	reverseMask string
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

	if len(t.mask) > MaxMaskLength {
		return nil, errors.New(fmt.Sprintf("Mask cannot be more than %d characters", MaxMaskLength))
	}

	t.reverseMask = stringReverse(t.mask)

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

// Utility for easing the template mask reversal
func stringReverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Generates a noid suffix for a given value of the noid sequence
func (template *Template) calculateSuffix(sequenceValue int64) string {
	var base int64
	var noidContainer [MaxMaskLength + 1]rune
	var i int
	var char rune

	// First, go through the mask in reverse, treating each mask character as a
	// base for the sequenceValue to convert to a noid character
	for i, char = range template.reverseMask {
		base = 10
		if char == 'e' {
			base = 29
		}
		templateChar := rune(ExtendedDigits[sequenceValue % int64(base)])
		noidContainer[MaxMaskLength - i] = templateChar
		sequenceValue = sequenceValue / base
	}

	// If sequenceValue wasn't completely used, and this isn't an "unlimited"
	// template, we can't mint a noid
	if sequenceValue > 0 && template.ordering != SequentialUnlimited {
		panic("sequenceValue out of range for template")
	}

	return string(noidContainer[MaxMaskLength-i:MaxMaskLength+1])
}
