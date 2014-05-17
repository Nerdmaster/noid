package noid

import(
	"errors"
)

type Minter struct {
	template *Template
	generator *SuffixGenerator
}

func NewMinter(template string, startSequence uint64) (*Minter, error) {
	t, err := NewTemplate(template)
	if err != nil {
		return nil, err
	}
	g := NewSuffixGenerator(t, startSequence)
	if g.totalBits > 64 {
		return nil, errors.New("Template range is too big!  Try a shorter template mask string.")
	}
	minter := &Minter{template: t, generator: g}

	return minter, nil
}

func (minter *Minter) Mint() string {
	noidSuffix := minter.generator.ToString()
	minter.generator.NextSequence()

	if minter.template.prefix == "" {
		return noidSuffix
	}

	return minter.template.prefix + "." + noidSuffix
}
