package noid

import "testing"

func assertEqualS(expected, actual string, message string, t *testing.T) {
	if (expected != actual) {
		t.Errorf("Expected %#v, but got %#v - %s", expected, actual, message)
	}
}

func TestToStringIsIdempotent(t *testing.T) {
	str := "foo.seedee"
	template, _ := NewTemplate(str)
	suffixGen := NewSuffixGenerator(template, 0)
	first := suffixGen.ToString()
	second := suffixGen.ToString()
	assertEqualS(first, second, "ToString shouldn't change suffixGen", t)
}

func TestSuffixGeneration_seedee(t *testing.T) {
	str := "foo.seedee"
	template, _ := NewTemplate(str)
	suffixGen := NewSuffixGenerator(template, 0)

	assertEqualS("00000", suffixGen.ToString(), "foo.seedee: 0", t)

	suffixGen.sequenceValue = 1
	assertEqualS("00001", suffixGen.ToString(), "foo.seedee: 1", t)

	suffixGen.sequenceValue = 1000
	assertEqualS("0015g", suffixGen.ToString(), "foo.seedee: 1000", t)

	suffixGen.sequenceValue = 100000
	assertEqualS("0c8w8", suffixGen.ToString(), "foo.seedee: 100000", t)

	// TODO: Test overflow
}

func TestSuffixGeneration_zdd(t *testing.T) {
	str := "foo.zdd"
	template, _ := NewTemplate(str)
	suffixGen := NewSuffixGenerator(template, 0)

	assertEqualS("00", suffixGen.ToString(), "foo.zdd: 0", t)

	suffixGen.sequenceValue = 1
	assertEqualS("01", suffixGen.ToString(), "foo.zdd: 1", t)

	suffixGen.sequenceValue = 1000
	assertEqualS("1000", suffixGen.ToString(), "foo.zdd: 1000", t)

	suffixGen.sequenceValue = 100000
	assertEqualS("100000", suffixGen.ToString(), "foo.zdd: 100000", t)
}
