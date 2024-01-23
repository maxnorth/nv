package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/maxnorth/nv/resolver"
	"github.com/spf13/cobra"
)

func RunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			// cmd can be 'run' or the 'nv' root - hard to
			// tell if this will run into trouble or not
			iRunCmd := cmd.ArgsLenAtDash()
			if iRunCmd == -1 {
				return errors.New("missing -- to begin command")
			}

			if iRunCmd != 0 {
				return errors.New("positional arguments are not allowed before --")
			}

			if len(args) == 0 {
				return errors.New("command was not provided after --")
			}

			r := resolver.Load()

			_, err := r.ResolveEnvironment()
			if err != nil {
				return err
			}

			execCommand(args)

			return nil
		},
	}

	return runCmd
}

func execCommand(commandArgs []string) {
	command, _ := commandArgs[0], []string{}
	fname, err := exec.LookPath(command)
	if err == nil {
		fname, err = filepath.Abs(fname)
	}
	if err != nil {
		log.Fatal(err)
	}
	err = syscall.Exec(fname, commandArgs, os.Environ())
	if err != nil {
		fmt.Printf("failed to run command: %s\n", err)
		os.Exit(1)
	}
}
