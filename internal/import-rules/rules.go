package importrules

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/payfazz/go-errors/v2"
	"gopkg.in/yaml.v3"
)

type rules []rule

type rule struct {
	ImportPath string   `yaml:"path"`
	Allowed    []string `yaml:"allowed"`
}

func (rules rules) normalize(mod string) {
	for i := range rules {
		rules[i].normalize(mod)
	}
}

func (rules rules) isValid(path, importing string) bool {
	for _, r := range rules {
		if r.isValid(path, importing) {
			return true
		}
	}
	return false
}

func (r *rule) normalize(mod string) {
	r.ImportPath = normalizeImportPath(mod, r.ImportPath)
	for i := range r.Allowed {
		r.Allowed[i] = normalizeImportPath(mod, r.Allowed[i])
	}
}

func (r *rule) isValid(path, importing string) bool {
	if !importPathMatch(path, r.ImportPath) {
		return false
	}
	for _, p := range r.Allowed {
		if importPathMatch(importing, p) {
			return true
		}
	}
	return false
}

func normalizeImportPath(mod, path string) string {
	if strings.HasPrefix(path, "+/") {
		mod = strings.TrimSuffix(mod, "/")
		return strings.TrimSuffix(mod+strings.TrimPrefix(path, "+"), "/")
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

func readRules(ctx context.Context, root string) (data []byte, err error) {
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

func loadRules(ctx context.Context, root, mod string) (rules, error) {
	data, err := readRules(ctx, root)
	if err != nil {
		return nil, err
	}
	var r rules
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, errors.Trace(err)
	}
	r.normalize(mod)
	return r, nil
}
