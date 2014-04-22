// Type and behaviors for generating a noid's suffix based on a template and an index

package noid

const Digits = "0123456789"
const ExtendedDigits = "0123456789bcdfghjkmnpqrstvwxz"

type SuffixContainer [MaxMaskLength]rune

type SuffixGenerator struct {
	sequenceValue int64
	index int
	minLength int
	suffix SuffixContainer
	reverseMaskBases []int64
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

func NewSuffixGenerator(template *Template, sequenceValue int64) *SuffixGenerator {
	nsg := &SuffixGenerator {sequenceValue: sequenceValue}
	nsg.ordering = template.ordering
	nsg.minLength = len(template.mask)

	reverseMask := stringReverseRunes(template.mask)
	nsg.reverseMaskBases = make([]int64, nsg.minLength)
	for i, char := range(reverseMask) {
		nsg.reverseMaskBases[i] = baseForMaskCharacter(char)
	}

	return nsg
}

// Returns the noid suffix for the given suffix generator - uses value, not
// pointer, to avoid altering the internal data
func (nsg SuffixGenerator) ToString() string {
	for nsg.sequenceValue > 0 || nsg.index < nsg.minLength {
		nsg.addCharacter()
	}

	return nsg.suffix.toString(nsg.index)
}

// Based on mask, ordering, and nsg state, prepends the next noid suffix char
func (nsg *SuffixGenerator) addCharacter() {
	base := nsg.reverseMaskBases[0]
	if len(nsg.reverseMaskBases) > 1 {
		nsg.reverseMaskBases = nsg.reverseMaskBases[1:]
	}

	// TODO: Re-add overflow detection

	val := nsg.sequenceValue % base

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
func baseForMaskCharacter(char rune) int64 {
	if char == 'd' {
		return 10
	}

	return 29
}
