package importrules

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/payfazz/go-errors/v2"
)

func Main(ctx context.Context) error {
	dir, mod, err := getGoMod(ctx)
	if err != nil {
		return err
	}

	rules, err := loadRules(dir, mod)
	if err != nil {
		return err
	}

	allImports, err := getAllImports(ctx, mod)
	if err != nil {
		return err
	}

	var paths []string
	for path := range allImports {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	isValid := true
	for _, path := range paths {
		imports := allImports[path]
		for _, imp := range imports {
			if !rules.isValid(path, imp) {
				fmt.Printf("import is not allowed: %s -> %s\n", path, imp)
				isValid = false
			}
		}
	}

	if !isValid {
		return errors.New("import rules validation error")
	}

	return nil
}

func getGoMod(ctx context.Context) (string, string, error) {
	stdout, err := command(ctx, "go", "list", "-m", "-f", `{{printf "%s\n%s" .Dir .Path}}`)
	if err != nil {
		return "", "", err
	}

	var mod, path string
	if s := strings.Split(stdout, "\n"); len(s) >= 2 {
		mod = s[0]
		path = s[1]
	}

	if mod == "" {
		return "", "", errors.New("no go.mod found")
	}

	mod = strings.TrimSuffix(mod, "/")

	return mod, path, nil
}

func getAllImports(ctx context.Context, mod string) (map[string][]string, error) {
	stdout, err := command(ctx, "go", "list", "-f", `{{printf "%s\n" .ImportPath }}{{range .Imports}}{{printf "%s\n" .}}{{end}}{{printf "=====\n"}}`, mod+"/...")
	if err != nil {
		return nil, err
	}
	ret := make(map[string][]string)

	list := strings.Split(stdout, "=====\n")
	for _, p := range list {
		if p == "" {
			continue
		}
		imports := strings.Split(p, "\n")
		if len(imports) > 0 {
			path := imports[0]
			for _, imp := range imports[1:] {
				if imp == "" {
					continue
				}
				ret[path] = append(ret[path], imp)
			}
		}
	}

	return ret, nil
}
