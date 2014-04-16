package noid

import "testing"

func assertEqualS(expected, actual string, message string, t *testing.T) {
	if (expected != actual) {
		t.Errorf("Expected %#v, but got %#v - %s", expected, actual, message)
	}
}

func assertTemplateAttributeS(templateString, attribute, expected, actual string, t *testing.T) {
	if (expected != actual) {
		t.Errorf("Expected %#v to have %s %#v, but got %#v", templateString, attribute, expected, actual)
	}
}

func assertTemplateAttributeO(templateString string, attribute string, expected Ordering, actual Ordering, t *testing.T) {
	if (expected != actual) {
		t.Errorf("Expected %#v to have %s %#v, but got %#v", templateString, attribute, expected, actual)
	}
}

func assertTemplateAttributeB(templateString string, attribute string, expected bool, actual bool, t *testing.T) {
	if (expected != actual) {
		t.Errorf("Expected %#v to have %s %#v, but got %#v", templateString, attribute, expected, actual)
	}
}

func TestTemplatesWithPrefix(t *testing.T) {
	str := "prefix.reedeek"
	template, _ := NewTemplate(str)

	assertTemplateAttributeS(str, "prefix", "prefix", template.prefix, t)
	assertTemplateAttributeS(str, "mask", "eedee", template.mask, t)
	assertTemplateAttributeO(str, "ordering", Random, template.ordering, t)
	assertTemplateAttributeB(str, "hasCheckDigit", true, template.hasCheckDigit, t)
}

func TestTemplatesWithoutPrefix(t *testing.T) {
	str := "reedeek"
	template, _ := NewTemplate(str)

	assertTemplateAttributeS(str, "prefix", "", template.prefix, t)
	assertTemplateAttributeS(str, "mask", "eedee", template.mask, t)
	assertTemplateAttributeO(str, "ordering", Random, template.ordering, t)
	assertTemplateAttributeB(str, "hasCheckDigit", true, template.hasCheckDigit, t)
}

func TestLongTemplates(t *testing.T) {
	str := "prefix.reedeedeedeedeedeedeedeedeek"
	template, err := NewTemplate(str)
	if template != nil {
		t.Errorf("Expected %#v to be invalid, but template was %#v", str, template)
	}

	if err == nil {
		t.Errorf("Expected %#v to be invalid, but err was nil", str)
	}
}

func TestBadOrder(t *testing.T) {
	str := "foo.xeedee"
	template, err := NewTemplate(str)

	if template != nil {
		t.Errorf("Expected %#v to be invalid, but template was non-nil", str)
	}

	if err == nil {
		t.Errorf("Expected %#v to be invalid, but err was nil", str)
	}
}

func TestMoreMoarMOAR(t *testing.T) {
	str := "foo.seedee"
	template, _ := NewTemplate(str)

	assertTemplateAttributeS(str, "prefix", "foo", template.prefix, t)
	assertTemplateAttributeS(str, "mask", "eedee", template.mask, t)
	assertTemplateAttributeO(str, "ordering", SequentialLimited, template.ordering, t)
	assertTemplateAttributeB(str, "hasCheckDigit", false, template.hasCheckDigit, t)

	str = "foo.zeedee"
	template, _ = NewTemplate(str)

	assertTemplateAttributeS(str, "prefix", "foo", template.prefix, t)
	assertTemplateAttributeS(str, "mask", "eedee", template.mask, t)
	assertTemplateAttributeO(str, "ordering", SequentialUnlimited, template.ordering, t)
	assertTemplateAttributeB(str, "hasCheckDigit", false, template.hasCheckDigit, t)
}

func TestNoidForIndex(t *testing.T) {
	str := "foo.seedee"
	template, _ := NewTemplate(str)
	assertEqualS("00000", template.calculateSuffix(0), "foo.seedee: 0", t)
	assertEqualS("00001", template.calculateSuffix(1), "foo.seedee: 1", t)
	assertEqualS("0015g", template.calculateSuffix(1000), "foo.seedee: 1000", t)
	assertEqualS("0c8w8", template.calculateSuffix(100000), "foo.seedee: 100000", t)

	str = "foo.zdd"
	template, _ = NewTemplate(str)
	assertEqualS("00", template.calculateSuffix(0), "foo.zdd: 0", t)
	assertEqualS("01", template.calculateSuffix(1), "foo.zdd: 1", t)
	assertEqualS("1000", template.calculateSuffix(1000), "foo.zdd: 1000", t)
	assertEqualS("100000", template.calculateSuffix(100000), "foo.zdd: 100000", t)
}
