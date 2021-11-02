package importrules

import "testing"

func TestNormalizeImportPath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		pkg      string
		expected string
	}{
		{"1a", "./", "mod"},
		{"1b", ".", "mod"},
		{"2", "aa/", "aa"},
		{"3", "./...", "mod/..."},
		{"4", "aa/...", "aa/..."},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if normalizeImportPath("mod", c.pkg) != c.expected {
				t.FailNow()
			}
		})
	}
}

func TestImportPathMatch(t *testing.T) {
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
		{"5", "bb", "aa/...", false},
		{"5", "x/y/z", "...", true},
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
	t.Parallel()
	r := rule{"./aa", []string{"bb", "cc/...", "./dd/..."}}
	r.normalize("mod")
	cases := []struct {
		name      string
		path      string
		importing string
		expected  bool
	}{
		{"1", "mod/aa", "bb", true},
		{"2", "mod/aa", "cc", true},
		{"3", "mod/aa", "cc/dd", true},
		{"4", "mod/bb", "bb", false},
		{"5", "mod/aa", "dd", false},
		{"6", "mod/aa", "mod/dd", true},
		{"7", "mod/aa", "mod/dd/ee", true},
		{"8", "mod/bb", "mod", false},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if r.isValid(c.path, c.importing) != c.expected {
				t.FailNow()
			}
		})
	}
}
