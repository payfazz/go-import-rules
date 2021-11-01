package importrules

import (
	"context"
	"fmt"
	"strings"

	"github.com/payfazz/go-errors/v2"
)

func Main(ctx context.Context) error {
	root, mod, err := getGoMod(ctx)
	if err != nil {
		return err
	}

	rules, err := loadRules(ctx, root, mod)
	if err != nil {
		return err
	}

	allImports, err := getAllImports(ctx, mod)
	if err != nil {
		return err
	}

	var paths []string
	for k := range allImports {
		paths = append(paths, k)
	}

	valid := true
	for _, path := range paths {
		imports := allImports[path]
		for _, i := range imports {
			if !rules.isValid(path, i) {
				fmt.Printf("import is not allowed: %s -> %s\n", path, i)
				valid = false
			}
		}
	}

	if !valid {
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
		pkg := strings.Split(p, "\n")
		if len(pkg) > 0 {
			h := pkg[0]
			for _, imp := range pkg[1:] {
				if imp == "" {
					continue
				}
				ret[h] = append(ret[h], imp)
			}
		}
	}

	return ret, nil
}
