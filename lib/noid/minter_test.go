package noid

import "testing"

func TestMinting(t *testing.T) {
	str := "foo.seedee"
	minter, _ := NewMinter(str, 0)
	assertEqualS("foo.00000", minter.Mint(), "foo.seedee first mint", t)

	for i := 0; i < 1000; i++ {
		minter.Mint()
	}
	assertEqualS("foo.000z9", minter.Mint(), "foo.seedee one-thousand-first mint", t)
}

func TestMintingWithNoPrefix(t *testing.T) {
	str := "zee"
	minter, _ := NewMinter(str, 0)
	assertEqualS("00", minter.Mint(), "zee first mint", t)

	for i := 0; i < 1000; i++ {
		minter.Mint()
	}
	assertEqualS("z9", minter.Mint(), "zee one-thousand-first mint", t)
}

func TestTemplateBitMaximums(t *testing.T) {
	str := "redededededededed"
	_, err := NewMinter(str, 0)
	if err != nil {
		t.Errorf("%s shouldn't have been too big (64 bits)!", str)
	}

	str = "reeeeeeeeeeeee"
	_, err = NewMinter(str, 0)
	if err == nil {
		t.Errorf("%s should have been too big (65 bits)!", str)
	}
}

func TestCheckdigitMagic(t *testing.T) {
	// It looks the same as above except for the check digit
	str := "foo.seedeek"
	minter, _ := NewMinter(str, 0)
	assertEqualS("foo.00000f", minter.Mint(), "foo.seedeek first mint", t)
	minter, _ = NewMinter(str, 1001)
	assertEqualS("foo.000z9r", minter.Mint(), "foo.seedeek one-thousand-first mint", t)

	// And the prefix matters for determining check digits... for some odd reason.
	str = "bar.seedeek"
	minter, _ = NewMinter(str, 0)
	assertEqualS("bar.000004", minter.Mint(), "bar.seedeek first mint", t)
	minter, _ = NewMinter(str, 1001)
	assertEqualS("bar.000z9d", minter.Mint(), "bar.seedeek one-thousand-first mint", t)
}
