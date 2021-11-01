package importrules

import "testing"

func TestNormalize(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		pkg      string
		expected string
	}{
		{"1", "+/", "mod"},
		{"2", "aa/", "aa"},
		{"3", "+/...", "mod/..."},
		{"4", "aa/...", "aa/..."},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if normalizeImportPath("mod", c.pkg) != c.expected {
				t.FailNow()
			}
			if normalizeImportPath("mod/", c.pkg) != c.expected {
				t.FailNow()
			}
		})
	}
}

func TestPkgMatch(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		pkg      string
		pat      string
		expected bool
	}{
		{"1", "aa", "aa", true},
		{"2", "aa", "aa/...", true},
		{"3", "aa/bb", "aa", false},
		{"4", "aa/bb", "aa/...", true},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if importPathMatch(c.pkg, c.pat) != c.expected {
				t.FailNow()
			}
		})
	}
}

func TestRuleIsValid(t *testing.T) {
	r := rule{"+/aa", []string{"bb", "cc/..."}}
	r.normalize("mod")
	if !r.isValid("mod/aa", "bb") {
		t.FailNow()
	}
	if !r.isValid("mod/aa", "cc") {
		t.FailNow()
	}
	if !r.isValid("mod/aa", "cc/dd") {
		t.FailNow()
	}
	if r.isValid("mod/bb", "bb") {
		t.FailNow()
	}
	if !r.isValid("mod/aa", "dd") {
		t.FailNow()
	}
}
