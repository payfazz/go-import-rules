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
	errstr := err.Error()
	if os.Getenv("GO_IMPORT_RULES_DEBUG") == "1" {
		errstr = errors.Format(err)
	}
	fmt.Fprintln(os.Stderr, errstr)
	return 1
}
