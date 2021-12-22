package importrules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/payfazz/go-errors/v2"
	"gopkg.in/yaml.v3"
)

type rules []rule

type rule struct {
	Path  string   `yaml:"path"`
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

func (rs rules) normalize(mod string) {
	for i := range rs {
		rs[i].normalize(mod)
	}
}

func (rs rules) isValid(pkg, imp string) bool {
	def := false
	for i := range rs {
		switch rs[i].decide(pkg, imp) {
		case undecided:
		case allowed:
			def = true
		case denied:
			return false
		default:
			panic("invalid switch")
		}
	}
	return def
}

func (r *rule) normalize(mod string) {
	r.Path = normalizeImportPath(mod, r.Path)
	for i := range r.Allow {
		r.Allow[i] = normalizeImportPath(mod, r.Allow[i])
	}
}

type decission int

const (
	undecided decission = iota
	allowed
	denied
)

func (r *rule) decide(pkg, imp string) decission {
	if !importPathMatch(pkg, r.Path) {
		return undecided
	}
	for _, p := range r.Deny {
		if importPathMatch(imp, p) {
			return denied
		}
	}
	for _, p := range r.Allow {
		if importPathMatch(imp, p) {
			return allowed
		}
	}
	return undecided
}

func normalizeImportPath(mod, path string) string {
	if path == "." {
		return mod
	}
	if strings.HasPrefix(path, "./") {
		path = mod + path[1:]
	}
	return strings.TrimSuffix(path, "/")
}

func importPathMatch(pkg, pattern string) bool {
	if pattern == "..." {
		return true
	}
	if !strings.HasSuffix(pattern, "/...") {
		return pkg == pattern
	}
	return strings.HasPrefix(pkg, pattern[:len(pattern)-4])
}

func readRules(root string) (data []byte, err error) {
	data, err = os.ReadFile(filepath.Join(root, "import-rules.yaml"))
	if err == nil {
		return
	}
	data, err = os.ReadFile(filepath.Join(root, "import-rules.yml"))
	if err == nil {
		return
	}
	return nil, errors.New(`cannot read "import-rules.yaml" or "import-rules.yml"`)
}

func loadRules(dir, mod string) (rules, error) {
	data, err := readRules(dir)
	if err != nil {
		return nil, nil
	}
	var r rules
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, errors.Trace(err)
	}
	r.normalize(mod)
	return r, nil
}
