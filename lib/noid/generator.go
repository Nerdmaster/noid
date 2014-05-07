// Type and behaviors for generating a noid's suffix based on a template and an index

package noid

import (
	"math"
	"errors"
)

// "d" is basically an octal value - this means each "d" uses precisely 3 bits
const Digits = "01234567"
const DigitBase = 8

// "e" is always a range of 32 characters - 5 bits
const ExtendedDigits = "0123456789abcdfghjkmnpqrstuvwxyz"
const ExtendedDigitBase = 32

type SuffixContainer [MaxMaskLength]rune

type SuffixGenerator struct {
	sequenceValue uint64
	maxSequence uint64
	index int
	minLength int
	suffix SuffixContainer
	reverseMaskBases []uint64
	ordering Ordering
	seed uint64
}

// Utility for easing the template mask reversal
func stringReverseRunes(s string) []rune {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return runes
}

func NewSuffixGenerator(template *Template, sequenceValue uint64) *SuffixGenerator {
	nsg := &SuffixGenerator {sequenceValue: sequenceValue}
	nsg.ordering = template.ordering
	nsg.minLength = len(template.mask)

	reverseMask := stringReverseRunes(template.mask)
	nsg.reverseMaskBases = make([]uint64, nsg.minLength)
	for i, char := range(reverseMask) {
		nsg.reverseMaskBases[i] = baseForMaskCharacter(char)
	}

	if nsg.ordering == SequentialUnlimited {
		nsg.maxSequence = math.MaxUint64
	} else {
		nsg.computeMaxSequenceValue()
	}

	// I dunno the right approach here - we want to be able to mint "random"
	// noids, but keep them predictable for a given template.
	if nsg.ordering == Random {
		nsg.seed = nsg.maxSequence / 2
	}

	if nsg.sequenceValue > nsg.maxSequence {
		return nil
	}

	return nsg
}

// Modifies seed, allowing for more randomness if an application needs this
func (nsg *SuffixGenerator) Seed(newSeed uint64) {
	nsg.seed = newSeed
}

// Computes the maximum sequence value by examining the bases (from the reverse
// mask array) for each character.  This is useful for both types of limited
// sequences to determine before computation if the sequence value is too high.
func (nsg *SuffixGenerator) computeMaxSequenceValue() {
	var multiplier uint64 = 1
	nsg.maxSequence = 0

	for _, base := range(nsg.reverseMaskBases) {
		nsg.maxSequence += multiplier * (base - 1)
		multiplier *= base
	}
}

// Returns the noid suffix for the given suffix generator - uses value, not
// pointer, to avoid altering the internal data
func (nsg SuffixGenerator) ToString() string {
	for nsg.sequenceValue > 0 || nsg.index < nsg.minLength {
		nsg.addCharacter()
	}

	return nsg.suffix.toString(nsg.index)
}

func (nsg *SuffixGenerator) NextSequence() error {
	if nsg.sequenceValue == nsg.maxSequence {
		return errors.New("Overflow trying to get next sequence")
	}

	nsg.sequenceValue++
	return nil
}

// Based on mask, ordering, and nsg state, prepends the next noid suffix char
func (nsg *SuffixGenerator) addCharacter() {
	base := nsg.reverseMaskBases[0]
	if len(nsg.reverseMaskBases) > 1 {
		nsg.reverseMaskBases = nsg.reverseMaskBases[1:]
	}

	val := nsg.sequenceValue

	if nsg.ordering == Random {
		val += nsg.seed
		nsg.seed -= 1
	}

	val %= base

	templateChar := rune(ExtendedDigits[val])
	nsg.suffix[MaxMaskLength - 1 - nsg.index] = templateChar
	nsg.sequenceValue /= base
	nsg.index++
}

func (nsc *SuffixContainer) toString(length int) string {
	return string(nsc[MaxMaskLength - length:MaxMaskLength])
}

// Uses hard-coded values 10 and 29 to quickly return the base a given
// character will be using
func baseForMaskCharacter(char rune) uint64 {
	if char == 'd' {
		return DigitBase
	}

	return ExtendedDigitBase
}
