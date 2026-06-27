package tr_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-tr"
	"github.com/gloo-foo/testable"
)

func ExampleTr_squeeze() {
	// echo "hello    world" | tr -s ' '
	output, _ := testable.Test(command.Tr(" ", "", command.TrSqueeze), "hello    world\n")
	fmt.Print(output)
	// Output:
	// hello world
}
