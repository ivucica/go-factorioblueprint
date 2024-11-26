# factorioblueprint

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/badc0de.net/pkg/factorioblueprint)

Package factorioblueprint is bringing basic tools to interact with Factorio blueprints to Go.

An experiment in how much can be done via code generation etc. Go package
names and "public" interfaces are both unstable at this point.

Using an older version of some modules so that Go 1.19 can be used (shipping
in Debian etc).

## Usage

Obtain and install the binary for the basic reader:

```go
$ go install badc0de.net/pkg/factorioblueprint/cmd/blueprintread@latest
```

which can then be printed out with the default prettyprint:

```go
$ blueprintread -file read_blueprint/simple.txt
```

or a format such as raw_json or yaml:

```go
$ blueprintread -fmt=raw_json -file read_blueprint/simple.txt
$ blueprintread -fmt=yaml -file read_blueprint/simple.txt
```

## Sub Packages

* [asciiart_blueprint](./asciiart_blueprint): Package asciiart_blueprint takes a blueprint schema and draws ASCII art for it.

* [read_blueprint](./read_blueprint)

* [write_blueprint](./write_blueprint)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
