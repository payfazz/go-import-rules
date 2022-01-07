package rules

import (
	"sort"
	"strings"

	"github.com/payfazz/go-errors/v2"
	"gopkg.in/yaml.v3"
)

type Rules []rulesItem

type rulesItem struct {
	pattern string
	allow   []string
	deny    []string
}

func (rs Rules) normalize(mod string) {
	for i := range rs {
		rs[i].normalize(mod)
	}
	sort.SliceStable(rs, func(i, j int) bool {
		return strings.Count(rs[i].pattern, "/") > strings.Count(rs[j].pattern, "/")
	})
}

func (rs Rules) IsAllowed(pkg, imp string) bool {
	for i := range rs {
		switch rs[i].decide(pkg, imp) {
		case undecided:
		case allowed:
			return true
		case denied:
			return false
		default:
			panic("invalid switch")
		}
	}
	return false
}

func (r *rulesItem) normalize(mod string) {
	r.pattern = normalizeImportPath(mod, r.pattern)
	for i := range r.allow {
		r.allow[i] = normalizeImportPath(mod, r.allow[i])
	}
	for i := range r.deny {
		r.deny[i] = normalizeImportPath(mod, r.deny[i])
	}
}

type decission int

const (
	undecided decission = iota
	allowed
	denied
)

func (r *rulesItem) decide(pkg, imp string) decission {
	if !importPathMatch(pkg, r.pattern) {
		return undecided
	}
	for _, p := range r.deny {
		if importPathMatch(imp, p) {
			return denied
		}
	}
	for _, p := range r.allow {
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

func ParseYAML(mod string, data []byte) (Rules, error) {
	type parsedItem struct {
		Path  string   `yaml:"path"`
		Allow []string `yaml:"allow"`
		Deny  []string `yaml:"deny"`
	}
	var parsed []parsedItem
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, errors.Trace(err)
	}
	var r Rules
	for _, item := range parsed {
		r = append(r, rulesItem{
			pattern: item.Path,
			allow:   item.Allow,
			deny:    item.Deny,
		})
	}
	r.normalize(mod)
	return r, nil
}
