// Package factorioblueprint is bringing basic tools to interact with Factorio blueprints to Go.
//
// An experiment in how much can be done via code generation etc. Go package
// names and "public" interfaces are both unstable at this point.
//
// Using an older version of some modules so that Go 1.19 can be used (shipping
// in Debian etc).
//
// Usage
//
//     $ go install badc0de.net/pkg/factorioblueprint/cmd/blueprintread@latest
//
// which can then be printed out with the default prettyprint:
//
//     $ blueprintread -file read_blueprint/simple.txt
//
// or a format such as raw_json or yaml:
//
//     $ blueprintread -fmt=raw_json -file read_blueprint/simple.txt
//     $ blueprintread -fmt=yaml -file read_blueprint/simple.txt
//
package factorioblueprint // badc0de.net/pkg/factorioblueprint
