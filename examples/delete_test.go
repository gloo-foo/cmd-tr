package tr_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-tr"
	"github.com/gloo-foo/testable"
)

func ExampleTr_delete() {
	// echo "hello123world" | tr -d '0-9'
	output, _ := testable.Test(command.Tr("0-9", "", command.TrDelete), "hello123world\n")
	fmt.Print(output)
	// Output:
	// helloworld
}
