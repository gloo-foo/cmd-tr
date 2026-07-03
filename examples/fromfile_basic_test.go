package tr_test

import (
	"fmt"
	"os"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/run"

	command "github.com/gloo-foo/cmd-tr"
)

// This example demonstrates reading from a file instead of inline input.
func ExampleTr_fromFile_basic() {
	// cat testdata/text.txt | tr 'a-z' 'A-Z'
	data, err := os.ReadFile("testdata/text.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	output, _ := testable.Test(command.Tr("a-z", "A-Z"), run.Input(string(data)))
	fmt.Print(output)
	// Output:
	// HELLO
}
