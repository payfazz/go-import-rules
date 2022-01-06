package importrules

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/payfazz/go-errors/v2"
)

// TODO(win): do parsing without calling go binary

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

	isValid := true
	for pkg, imports := range allImports {
		for _, imp := range imports {
			if !rules.isAllowed(pkg, imp) {
				fmt.Printf("import is not allowed: %s -> %s\n", pkg, imp)
				isValid = false
			}
		}
	}

	if !isValid {
		os.Exit(2)
	}

	return nil
}

func getGoMod(ctx context.Context) (dir, mod string, err error) {
	stdout, err := command(ctx, "go", "list", "-m", "-f", `{{printf "%s\n%s" .Dir .Path}}`)
	if err != nil {
		return "", "", err
	}

	if s := strings.Split(stdout, "\n"); len(s) >= 2 {
		dir = s[0]
		mod = s[1]
	}

	if dir == "" {
		return "", "", errors.New("no go.mod found")
	}

	dir = strings.TrimSuffix(dir, string(os.PathSeparator))
	mod = strings.TrimSuffix(mod, "/")

	return dir, mod, nil
}

func getAllImports(ctx context.Context, mod string) (map[string][]string, error) {
	stdout, err := command(ctx, "go", "list", "-f", `{{printf ">%s\n" .ImportPath}}{{range .Imports}}{{printf "+%s\n" .}}{{end}}`, mod+"/...")
	if err != nil {
		return nil, err
	}

	imports := make(map[string][]string)

	lines := strings.Split(stdout, "\n")
	var curPkg string
	for _, line := range lines {
		if line == "" {
			continue
		}
		mark, pkg := line[0], line[1:]
		if mark == '>' {
			curPkg = pkg
			continue
		}
		if mark == '+' {
			if curPkg == "" || pkg == "" {
				continue
			}
			imports[curPkg] = append(imports[curPkg], pkg)
			continue
		}
		panic("invalid mark")
	}

	return imports, nil
}
