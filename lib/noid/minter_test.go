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
