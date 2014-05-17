// Type and behaviors for generating a noid's suffix based on a template and an index

package noid

import (
	"math"
	"errors"
)

// "d" is basically an octal value - this means each "d" uses precisely 3 bits
const Digits = "01234567"
const DigitBits = 3

// "e" is always a range of 32 characters - 5 bits
const ExtendedDigits = "0123456789abcdfghjkmnpqrstuvwxyz"
const ExtendedDigitBits = 5

// Maximum possible mask within 64 bits (if using "d" for all characters)
const MaxMaskLength = 21

type SuffixContainer [MaxMaskLength]rune

type SuffixGenerator struct {
	sequenceValue uint64
	maxSequence uint64
	index int
	minLength int
	suffix SuffixContainer
	reverseMaskBits []byte
	totalBits byte
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
	nsg.reverseMaskBits = make([]byte, nsg.minLength)
	for i, char := range(reverseMask) {
		nsg.reverseMaskBits[i] = bitsForMaskCharacter(char)
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

// Computes the maximum sequence value by examining the bit size (from the reverse
// mask array) for each character.  This is useful for both types of limited
// sequences to determine before computation if the sequence value is too high.
func (nsg *SuffixGenerator) computeMaxSequenceValue() {
	nsg.totalBits = 0

	for _, bits := range(nsg.reverseMaskBits) {
		nsg.totalBits += bits
	}

	nsg.maxSequence = (1 << nsg.totalBits) - 1
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
	bits := nsg.reverseMaskBits[0]
	if len(nsg.reverseMaskBits) > 1 {
		nsg.reverseMaskBits = nsg.reverseMaskBits[1:]
	}

	val := nsg.sequenceValue

	if nsg.ordering == Random {
		val += nsg.seed
		nsg.seed -= 1
	}

	val = val & ((1 << bits) - 1)

	templateChar := rune(ExtendedDigits[val])
	nsg.suffix[MaxMaskLength - 1 - nsg.index] = templateChar
	nsg.sequenceValue >>= bits
	nsg.index++
}

func (nsc *SuffixContainer) toString(length int) string {
	return string(nsc[MaxMaskLength - length:MaxMaskLength])
}

// Returns bit constants for a given mask character
func bitsForMaskCharacter(char rune) byte {
	if char == 'd' {
		return DigitBits
	}

	return ExtendedDigitBits
}
