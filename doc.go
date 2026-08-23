// Package factorioblueprint is bringing basic tools to interact with Factorio blueprints to Go.
//
// An experiment in how much can be done via code generation etc. Go package
// names and "public" interfaces are both unstable at this point.
//
// Using an older version of some modules so that Go 1.19 can be used (shipping
// in Debian etc).
//
// NOTE: Switched to 1.21 solely so to avoid needing to triage why 1.21 is being
// pulled in.
//
// Usage
//
// Obtain and install the binary for the basic reader:
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
// README BUG: Above should be bash codeblocks, but goreadme always uses Go.
package factorioblueprint // badc0de.net/pkg/factorioblueprint
