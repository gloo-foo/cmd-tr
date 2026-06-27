// Package alias provides short names for tr command flags.
package alias

import command "github.com/gloo-foo/cmd-tr"

// Tr is the tr command constructor.
var Tr = command.Tr

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
