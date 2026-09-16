// Package alias provides short names for tr command flags.
package alias

import (
	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-tr"
)

// Tr builds the tr command: from is set1, to is set2; ranges like "a-z" expand.
func Tr(from, to command.TrSet, opts ...any) gloo.Command[[]byte, []byte] {
	return command.Tr(from, to, opts...)
}

// Set names a tr character set argument (set1/set2).
type Set = command.TrSet

// -d flag: delete characters
const Delete = command.TrDelete

// default: translate
const NoDelete = command.TrNoDelete

// -s flag: squeeze repeats
const Squeeze = command.TrSqueeze

// default: don't squeeze
const NoSqueeze = command.TrNoSqueeze

// -c flag: complement set
const Complement = command.TrComplement

// default: use set as-is
const NoComplement = command.TrNoComplement
