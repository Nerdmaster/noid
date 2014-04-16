package noid

type Minter struct {
	index int64
	template *Template
}

func NewMinter(template string) (*Minter, error) {
	t, err := NewTemplate(template)
	if err != nil {
		return nil, err
	}
	minter := &Minter{index: 0, template: t}

	return minter, nil
}

func (minter *Minter) Mint() string {
	noidSuffix := minter.template.calculateSuffix(minter.index)
	minter.index++

	if minter.template.prefix == "" {
		return noidSuffix
	}

	return minter.template.prefix + "." + noidSuffix
}
