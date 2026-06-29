package tr_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-tr"
)

func ExampleTr_delete() {
	// echo "hello123world" | tr -d '0-9'
	output, _ := testable.Test(command.Tr("0-9", "", command.TrDelete), "hello123world\n")
	fmt.Print(output)
	// Output:
	// helloworld
}
