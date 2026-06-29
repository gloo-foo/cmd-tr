package tr_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-tr"
)

func ExampleTr_squeeze() {
	// echo "hello    world" | tr -s ' '
	output, _ := testable.Test(command.Tr(" ", "", command.TrSqueeze), "hello    world\n")
	fmt.Print(output)
	// Output:
	// hello world
}
