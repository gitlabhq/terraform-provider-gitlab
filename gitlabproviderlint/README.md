# gitlabproviderlint

The `gitlabproviderlint` tool is a code linting tool, specifically tailored for the Terraform GitLab Provider.

## Checks

### Documentation Checks

Check | Description 
--- | ---
GD01 | check for resources without a docs page

## Run

Analyze the `./gitlab` package:

```sh
$ go run ./gitlabproviderlint ./gitlab
```

## Development

This tool uses the standard [go/analysis](https://godoc.org/golang.org/x/tools/go/analysis) API for implementing custom static code analysis.

### Unit testing

```sh
$ go test ./gitlabproviderlint
```

### Adding an analyzer

1. Pick a unique name, using the `G` prefix to help avoid naming collisions with external analyzers. (i.e. "GD##"" for a docs-related analyzer.)
1. Add `NAME.go` (where `NAME` is your analyzer name) and implement [Analyzer](https://godoc.org/golang.org/x/tools/go/analysis#Analyzer).
1. Add `NAME_test.go` with a test using [analysistest.Run()](https://godoc.org/golang.org/x/tools/go/analysis/analysistest#Run).
1. Add a `testdata/src/NAME` directory with Go source files that implement passing and failing code based on [analysistest](https://godoc.org/golang.org/x/tools/go/analysis/analysistest) framework.
1. Since [analysistest](https://godoc.org/golang.org/x/tools/go/analysis/analysistest) does not support Go Modules currently, each analyzer that implements testing must add a symlink to the top level vendor directory in the Go package beneath `testdata/NAME/src`.
   - Example: `ln -s ../../../../vendor testdata/src/NAME`
   - Example: `ln -s ../../../../../vendor testdata/src/NAME/gitlab` for a nested package.
1. Add the analyzer to the `allAnalyzers` variable in `main.go`.
1. Add the analyzer to the table in this `README.md`.
