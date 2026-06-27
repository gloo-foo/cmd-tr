package tr_test

import (
	"fmt"
	"os"

	command "github.com/gloo-foo/cmd-tr"
	"github.com/gloo-foo/testable"
)

// This example demonstrates reading from a file instead of inline input.
func ExampleTr_fromFile_basic() {
	// cat testdata/text.txt | tr 'a-z' 'A-Z'
	data, err := os.ReadFile("testdata/text.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read testdata: %v\n", err)
		return
	}
	output, _ := testable.Test(command.Tr("a-z", "A-Z"), string(data))
	fmt.Print(output)
	// Output:
	// HELLO
}
