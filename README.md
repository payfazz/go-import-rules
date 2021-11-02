# go-import-rules

Utility to restrict which package is allowed to import another package.

This tool will read `import-rules.yaml` or `import-rules.yml` in the the folder that contains `go.mod`

## Installation

```sh
go install github.com/payfazz/go-import-rules/...@latest
```

## Usage

Create `import-rules.yaml`, this file contains array of rules.

Each rule have `path` (`string`) and `allow` (`[]string`) property.

Import path in rules have following convention:
- path start with `"./"` (or `"."`) mean path that spesified in `go.mod`.
- path end with `"/..."` will include subpath (for more info, run `go help packages`).

Then, run `go-import-rules`, if you got no error message, that mean you are good to go.
