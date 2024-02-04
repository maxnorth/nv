package cmd_test

import (
	"fmt"
	"testing"

	"github.com/maxnorth/nv/test"
)

func Test_RunCmd(t *testing.T) {
	exitCode, stdout, stderr := test.RunCLI(t, "basic", "--", "echo", "test6")
	fmt.Print(exitCode, stdout, stderr)
}
