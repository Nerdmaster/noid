package noid

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
	minter := &Minter{template: t, generator: g}

	return minter, nil
}

func (minter *Minter) Mint() string {
	noidSuffix := minter.generator.ToString()
	minter.generator.sequenceValue++

	if minter.template.prefix == "" {
		return noidSuffix
	}

	return minter.template.prefix + "." + noidSuffix
}
