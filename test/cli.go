package test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/maxnorth/nv/cmd"
)

func RunCLI(t *testing.T, testCase string, args ...string) (exitCode int, outstr string, errstr string) {
	// return runCLIRootCmd(t, testCase, args...)
	return runCLIBinary(t, testCase, args...)
}

func runCLIRootCmd(t *testing.T, testCase string, args ...string) (exitCode int, strout string, strerr string) {
	testCaseDir := path.Join(RootDir(), "test", "cases", testCase)
	os.Chdir(testCaseDir)

	cmd := cmd.NewRootCmd()
	bufin := bytes.NewBufferString("")
	bufout := bytes.NewBufferString("")
	buferr := bytes.NewBufferString("")

	cmd.SetIn(bufin)
	cmd.SetOut(bufout)
	cmd.SetErr(buferr)
	cmd.SetArgs(args)

	cmd.Execute()
	byteout, err := io.ReadAll(bufout)
	if err != nil {
		t.Fatal(err)
	}
	byteerr, err := io.ReadAll(buferr)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: accurate exitCode result
	return 0, string(byteout), string(byteerr)
}

func runCLIBinary(t *testing.T, testCase string, args ...string) (exitCode int, strout string, strerr string) {
	nvPath := path.Join(RootDir(), "dist", "nv")
	cmd := exec.Command(nvPath, args...)

	cmd.Dir = path.Join(RootDir(), "test", "cases", testCase)

	bufout := bytes.NewBufferString("")
	buferr := bytes.NewBufferString("")
	cmd.Stdout = bufout
	cmd.Stderr = buferr

	err := cmd.Run()

	byteout, err := io.ReadAll(bufout)
	if err != nil {
		t.Fatal(err)
	}
	byteerr, err := io.ReadAll(buferr)
	if err != nil {
		t.Fatal(err)
	}

	return cmd.ProcessState.ExitCode(), string(byteout), string(byteerr)
}
