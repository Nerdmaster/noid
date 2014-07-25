package noid

import (
	"testing"
	"math"
)

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
	assertEqualS("000z8", suffixGen.ToString(), "foo.seedee: 1000", t)

	suffixGen.sequenceValue = 100000
	assertEqualS("0c1p0", suffixGen.ToString(), "foo.seedee: 100000", t)

	// TODO: Test overflow
}

func TestSuffixGeneration_zdd(t *testing.T) {
	str := "foo.zdd"
	template, _ := NewTemplate(str)
	suffixGen := NewSuffixGenerator(template, 0)

	assertEqualS("00", suffixGen.ToString(), "foo.zdd: 0", t)

	suffixGen.sequenceValue = 1
	assertEqualS("01", suffixGen.ToString(), "foo.zdd: 1", t)

	// Verify octals real quick
	suffixGen.sequenceValue = 01000
	assertEqualS("1000", suffixGen.ToString(), "foo.zdd: 01000", t)

	// And verify octals again!
	suffixGen.sequenceValue = 0100000
	assertEqualS("100000", suffixGen.ToString(), "foo.zdd: 0100000", t)
}

func TestSuffixMaximum(t *testing.T) {
	var template *Template
	var g *SuffixGenerator

	template, _ = NewTemplate("sdd")
	g = NewSuffixGenerator(template, uint64(math.Exp2(3 + 3)) - 1)
	assertEqualUint64(g.sequenceValue, g.maxSequence, "Max sequence for sdd", t)
	assertEqualS("77", g.ToString(), "Suffix using maximum sequence", t)

	template, _ = NewTemplate("sedd")
	g = NewSuffixGenerator(template, uint64(math.Exp2(5 + 3 + 3)) - 1)
	assertEqualUint64(g.sequenceValue, g.maxSequence, "Max sequence for sedd", t)
	assertEqualS("z77", g.ToString(), "Suffix using maximum sequence", t)

	template, _ = NewTemplate("sdedd")
	g = NewSuffixGenerator(template, uint64(math.Exp2(3 + 5 + 3 + 3)) - 1)
	assertEqualUint64(g.sequenceValue, g.maxSequence, "Max sequence for sdedd", t)
	assertEqualS("7z77", g.ToString(), "Suffix using maximum sequence", t)

	template, _ = NewTemplate("zdd")
	g = NewSuffixGenerator(template, 0)
	assertEqualUint64(18446744073709551615, g.maxSequence, "Max sequence for zdd: uint64 max", t)
}

func TestNextSequence(t *testing.T) {
	var template *Template
	var g *SuffixGenerator

	template, _ = NewTemplate("sdd")
	g = NewSuffixGenerator(template, 0)
	g.NextSequence()
	assertEqualUint64(1, g.sequenceValue, "0 + 1 is 1!!!!", t)
}

func TestRandomNoids(t *testing.T) {
	var template *Template
	var g *SuffixGenerator

	template, _ = NewTemplate("reedee")
	g = NewSuffixGenerator(template, 0)
	assertEqualS("q67j4", g.ToString(), "reedee @ sequence 0", t)

	g.sequenceValue = g.maxSequence / 2
	assertEqualS("tt0fv", g.ToString(), "reedee @ middle sequence", t)

	g.sequenceValue = g.maxSequence
	assertEqualS("9t0fv", g.ToString(), "reedee @ last sequence", t)
}

func TestRandomNoidsDontRepeat(t *testing.T) {
	var template *Template
	var g *SuffixGenerator
	var err error

	template, _ = NewTemplate("reee")
	g = NewSuffixGenerator(template, 0)

	seen := make(map[string]bool)

	for err == nil {
		if seen[g.ToString()] {
			t.Errorf("ACK!!  %#v was seen before!", g.ToString())
		}

		seen[g.ToString()] = true
		err = g.NextSequence()
	}
}
