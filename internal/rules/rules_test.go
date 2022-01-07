package rules

import "testing"

func TestNormalizeImportPath(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		pkg      string
		expected string
	}{
		"1a": {"./", "mod"},
		"1b": {".", "mod"},
		"2":  {"aa/", "aa"},
		"3":  {"./...", "mod/..."},
		"4":  {"aa/...", "aa/..."},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if normalizeImportPath("mod", test.pkg) != test.expected {
				t.FailNow()
			}
		})
	}
}

func TestImportPathMatch(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		pkg      string
		pat      string
		expected bool
	}{
		"1": {"aa", "aa", true},
		"2": {"aa", "aa/...", true},
		"3": {"aa/bb", "aa", false},
		"4": {"aa/bb", "aa/...", true},
		"5": {"bb", "aa/...", false},
		"6": {"x/y/z", "...", true},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if importPathMatch(test.pkg, test.pat) != test.expected {
				t.FailNow()
			}
		})
	}
}

func TestRuleDecide(t *testing.T) {
	t.Parallel()
	r := rulesItem{"./aa", []string{"bb", "cc/...", "./dd/...", "xx"}, []string{"xx", "yy"}}
	r.normalize("mod")
	tests := map[string]struct {
		path      string
		importing string
		expected  decission
	}{
		"1":  {"mod/aa", "bb", allowed},
		"2":  {"mod/aa", "cc", allowed},
		"3":  {"mod/aa", "cc/dd", allowed},
		"4":  {"mod/bb", "bb", undecided},
		"5":  {"mod/aa", "dd", undecided},
		"6":  {"mod/aa", "mod/dd", allowed},
		"7":  {"mod/aa", "mod/dd/ee", allowed},
		"8":  {"mod/bb", "mod", undecided},
		"9":  {"mod/aa", "xx", denied},
		"10": {"mod/aa", "yy", denied},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if r.decide(test.path, test.importing) != test.expected {
				t.FailNow()
			}
		})
	}
}

func TestRulesDenySubPath(t *testing.T) {
	t.Parallel()
	rs := Rules{
		{"./aa/...", []string{"./aa/..."}, []string{}},
		{"./aa/bb/...", []string{}, []string{"./aa/cc/..."}},
	}
	rs.normalize("mod")
	tests := map[string]struct {
		path      string
		importing string
		allowed   bool
	}{
		"1": {"mod/aa/xyz", "mod/aa/ppp", true},
		"2": {"mod/aa/bb/xyz", "mod/aa/cc/ppp", false},
		"3": {"some", "other", false},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if rs.IsAllowed(test.path, test.importing) != test.allowed {
				t.FailNow()
			}
		})
	}
}

func TestRulesAllowAllExceptSome(t *testing.T) {
	t.Parallel()
	rs := Rules{
		{"...", []string{"..."}, []string{}},
		{"./aa/...", []string{"fmt", "strings"}, []string{}},
		{"./aa/...", []string{}, []string{"..."}},
	}
	rs.normalize("mod")
	tests := map[string]struct {
		path      string
		importing string
		allowed   bool
	}{
		"1": {"mod/some", "something", true},
		"2": {"mod/aa/some", "something", false},
		"3": {"mod/aa/some", "fmt", true},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if rs.IsAllowed(test.path, test.importing) != test.allowed {
				t.FailNow()
			}
		})
	}
}

func TestParseYAML(t *testing.T) {
	t.Parallel()
	data := `` +
		`- path: ./...` + "\n" +
		`  allow:` + "\n" +
		`    - ...` + "\n" +
		``

	_, err := ParseYAML("mod", []byte(data))
	if err != nil {
		t.FailNow()
	}
}
