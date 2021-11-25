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
}

func (rs rules) normalize(mod string) {
	for i := range rs {
		rs[i].normalize(mod)
	}
}

func (rs rules) isValid(path, importing string) bool {
	for i := range rs {
		if rs[i].isValid(path, importing) {
			return true
		}
	}
	return false
}

func (r *rule) normalize(mod string) {
	r.Path = normalizeImportPath(mod, r.Path)
	for i := range r.Allow {
		r.Allow[i] = normalizeImportPath(mod, r.Allow[i])
	}
}

func (r *rule) isValid(path, importing string) bool {
	if !importPathMatch(path, r.Path) {
		return false
	}
	for _, p := range r.Allow {
		if importPathMatch(importing, p) {
			return true
		}
	}
	return false
}

func normalizeImportPath(mod, path string) string {
	if path == "." {
		return mod
	}
	if strings.HasPrefix(path, "./") {
		path = mod + strings.TrimPrefix(path, ".")
	}
	return strings.TrimSuffix(path, "/")
}

func importPathMatch(path, pattern string) bool {
	if pattern == "..." {
		return true
	}
	if !strings.HasSuffix(pattern, "/...") {
		return path == pattern
	}
	return strings.HasPrefix(path, strings.TrimSuffix(pattern, "/..."))
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

func loadRules(root, mod string) (rules, error) {
	data, err := readRules(root)
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
