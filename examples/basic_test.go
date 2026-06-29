package tr_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-tr"
)

func ExampleTr_basic() {
	// echo "hello" | tr 'a-z' 'A-Z'
	output, _ := testable.Test(command.Tr("a-z", "A-Z"), "hello\n")
	fmt.Print(output)
	// Output:
	// HELLO
}
