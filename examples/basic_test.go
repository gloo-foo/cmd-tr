package tr_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-tr"
	"github.com/gloo-foo/testable"
)

func ExampleTr_basic() {
	// echo "hello" | tr 'a-z' 'A-Z'
	output, _ := testable.Test(command.Tr("a-z", "A-Z"), "hello\n")
	fmt.Print(output)
	// Output:
	// HELLO
}
