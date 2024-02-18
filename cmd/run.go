package cmd

import (
	"errors"
	"os/exec"

	"github.com/maxnorth/nv/internal/resolver"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
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

			r, err := resolver.Load(cmd.Flag("env").Value.String())
			if err != nil {
				return err
			}

			_, err = r.ResolveEnvironment()
			if err != nil {
				return err
			}

			err = execCommand(cmd, args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return runCmd
}

func execCommand(cobraCmd *cobra.Command, commandArgs []string) error {
	command, args := commandArgs[0], []string{}
	if len(commandArgs) > 1 {
		args = commandArgs[1:]
	}
	cmd := exec.Command(command, args...)
	cmd.Dir = initWorkDir
	cmd.Stdin = cobraCmd.InOrStdin()
	cmd.Stdout = cobraCmd.OutOrStdout()
	cmd.Stderr = cobraCmd.ErrOrStderr()
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
