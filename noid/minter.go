package noid

import (
	"errors"
	"strings"
)

type Minter struct {
	template  *Template
	generator *SuffixGenerator
}

func NewMinter(template string) (*Minter, error) {
	return NewSequencedMinter(template, 0)
}

func NewSequencedMinter(template string, startSequence uint64) (*Minter, error) {
	t, err := NewTemplate(template)
	if err != nil {
		return nil, err
	}
	g := NewSuffixGenerator(t, startSequence)
	if g == nil {
		return nil, errors.New("Minter sequence value too high")
	}
	if g.totalBits > 64 {
		return nil, errors.New("Template range is too big!  Try a shorter template mask string.")
	}
	minter := &Minter{template: t, generator: g}

	return minter, nil
}

func (m *Minter) Template() string {
	return m.template.String()
}

func (m *Minter) Sequence() uint64 {
	return m.generator.Sequence()
}

func (minter *Minter) Mint() string {
	result := minter.generator.ToString()
	minter.generator.NextSequence()

	if minter.template.Prefix != "" {
		result = minter.template.Prefix + "." + result
	}

	if minter.template.HasCheckDigit {
		result = result + string(computeCheckDigit(result))
	}

	return result
}

func computeCheckDigit(s string) rune {
	tally := 0
	runes := []rune(ExtendedDigits)
	for index, ch := range s {
		idx := strings.IndexRune(ExtendedDigits, ch)
		if idx == -1 {
			idx = 0
		}
		tally += idx * (1 + index)
	}
	return runes[tally%len(ExtendedDigits)]
}
