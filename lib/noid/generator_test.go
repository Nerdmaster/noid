package noid

import "testing"

func assertEqualUint64(expected, actual uint64, message string, t *testing.T) {
	if (expected != actual) {
		t.Errorf("Expected %d, but got %d - %s", expected, actual, message)
	}
}

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

func TestSuffixMaximum(t *testing.T) {
	var template *Template
	var g *SuffixGenerator

	template, _ = NewTemplate("sdd")
	g = NewSuffixGenerator(template, 99)
	assertEqualUint64(g.sequenceValue, g.maxSequence, "Max sequence for sdd", t)
	assertEqualS("99", g.ToString(), "Suffix using maximum sequence", t)

	template, _ = NewTemplate("sedd")
	g = NewSuffixGenerator(template, 2899)
	assertEqualUint64(g.sequenceValue, g.maxSequence, "Max sequence for sedd", t)
	assertEqualS("z99", g.ToString(), "Suffix using maximum sequence", t)

	template, _ = NewTemplate("sdedd")
	g = NewSuffixGenerator(template, 28999)
	assertEqualUint64(g.sequenceValue, g.maxSequence, "Max sequence for sdedd", t)
	assertEqualS("9z99", g.ToString(), "Suffix using maximum sequence", t)

	template, _ = NewTemplate("zdd")
	g = NewSuffixGenerator(template, 0)
	assertEqualUint64(18446744073709551615, g.maxSequence, "Max sequence for zdd: uint64 max", t)
}
