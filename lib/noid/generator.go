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

	if nsg.sequenceValue > nsg.maxSequence {
		return nil
	}

	return nsg
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

// Shuffles bits and xors stuff to map one sequence to another
//
// TODO: Cache bit pairs so we don't recompute this on every single iteration
func (nsg *SuffixGenerator) randomizeSequence() {
	var maxBit byte = nsg.totalBits - 1
	var bitIndex byte

	// 2/3 of the max value gives us a repeating "1010..." bit pattern, which is
	// a decent xor value to start with
	xor := nsg.maxSequence * 2 / 3

	// Create a changing seed based on our xor value for bit swapping "randomness"
	seed := xor

	// Temporary local var to ease code (and possibly avoid indirection)
	sval := nsg.sequenceValue ^ xor

	// Make sure the lowest bits are distributed a little - we always have at
	// least three bits, so this will never crash, though it won't necessarily be
	// all that useful, either.
	//
	// 3 bits:  0 and 2 			1 and 1
	// 5 bits:  0 and 3 			1 and 2
	// 20 bits: 0 and 18			1 and 9
	sval = bitSwap(sval, 0, maxBit - 1)
	sval = bitSwap(sval, 1, maxBit >> 1)

	for bitIndex = 3; bitIndex < maxBit; bitIndex++ {
		bit2 := seed % uint64(nsg.totalBits)
		sval = bitSwap(sval, bitIndex, byte(bit2))
		seed = seed >> 1
	}

	nsg.sequenceValue = sval
}

// Returns the noid suffix for the given suffix generator - uses value, not
// pointer, to avoid altering the internal data
func (nsg SuffixGenerator) ToString() string {
	if nsg.ordering == Random {
		nsg.randomizeSequence()
	}

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

	val := nsg.sequenceValue & ((1 << bits) - 1)

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
