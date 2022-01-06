package main

import (
	"fmt"
	"os"

	"github.com/payfazz/go-errors/v2"
	"github.com/payfazz/go-mainrun"

	importrules "github.com/payfazz/go-import-rules/internal/import-rules"
)

func main() {
	mainrun.OnError(onError)
	mainrun.Run(importrules.Main)
}

func onError(err error) int {
	errstr := "Error: " + err.Error()
	if os.Getenv("GO_IMPORT_RULES_STACKTRACE") == "1" {
		errstr = errors.FormatWithFilterPkgs(err, "main", "github.com/payfazz/go-import-rules")
	}
	fmt.Fprintln(os.Stderr, errstr)
	var detailedError interface{ ErrorDetail() string }
	if errors.As(err, &detailedError) {
		fmt.Fprintln(os.Stderr, detailedError.ErrorDetail())
	}
	return 1
}
