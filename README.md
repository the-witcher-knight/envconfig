# envconfig

A lightweight, zero-dependency Go package for loading configuration values directly from OS environment variables. This package provides custom tag-based validation to ensure that required environment variables are set and meet specified conditions.

![Alt](https://repobeats.axiom.co/api/embed/542629544c93a3a78e4a0b2fcd722f49217ce11c.svg "Repobeats analytics image")

## Features

- **Zero Dependencies**: No external libraries required.
- **Custom Tag-Based Validation**: Use struct tags to validate environment variables for required fields, format, and value constraints.
- **Simple and Lightweight**: Minimal overhead for maximum performance.

## Installation

```bash
go get -u github.com/the-witcher-knight/envconfig
```

## Unit test

```
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
        github.com/the-witcher-knight/envconfig/examples                coverage: 0.0% of statements
=== RUN   Test_Lookup
=== RUN   Test_Lookup/success
=== RUN   Test_Lookup/error_-_validation_fail
--- PASS: Test_Lookup (0.00s)
    --- PASS: Test_Lookup/success (0.00s)
    --- PASS: Test_Lookup/error_-_validation_fail (0.00s)
=== RUN   Test_AddValidator
--- PASS: Test_AddValidator (0.00s)
PASS
coverage: 82.1% of statements in ./...
ok      github.com/the-witcher-knight/envconfig 0.624s  coverage: 82.1% of statements in ./...

```

## Benchmark

```
go test -bench=.
goos: darwin
goarch: arm64
pkg: github.com/the-witcher-knight/envconfig
cpu: Apple M1 Pro
Benchmark_Lookup-10       230292              5234 ns/op            1283 B/op         56 allocs/op
PASS
ok      github.com/the-witcher-knight/envconfig 2.797s
```
