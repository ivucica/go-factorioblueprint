# factorioblueprint

Bringing basic tools to interact with Factorio blueprints to Go.

An experiment in how much can be done via code generation etc.

Go package names and "public" interfaces are unstable at this point.

Using an older version of some modules so that Go 1.19 can be used (shipping in
Debian etc).

## Usage

```bash
go install badc0de.net/pkg/factorioblueprint/cmd/blueprintread@latest
```

then default prettyprint:

```bash
blueprintread -file read_blueprint/simple.txt
```

or format such as raw_json or yaml:

```bash
blueprintread -fmt=raw_json -file read_blueprint/simple.txt
blueprintread -fmt=yaml -file read_blueprint/simple.txt
```
