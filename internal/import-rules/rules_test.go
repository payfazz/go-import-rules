package importrules

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
	r := rule{"./aa", []string{"bb", "cc/...", "./dd/...", "xx"}, []string{"xx", "yy"}}
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
