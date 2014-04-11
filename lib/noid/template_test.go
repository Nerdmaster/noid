package noid

import "testing"

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
		t.Errorf("Expected %#v to be invalid, but template was non-nil", str)
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
