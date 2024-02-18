package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initWorkDir string

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "nv",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			workDir, _ := os.Getwd()
			initWorkDir = workDir

			dir := cmd.Flag("dir").Value.String()
			if !strings.HasPrefix(dir, "/") {
				dir = filepath.Join(workDir, dir)
			}

			if err := os.Chdir(dir); err != nil {
				// TODO: other kinds of errors? permission issues?
				return fmt.Errorf("target directory not found: %s", dir)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// if -- not in args
			if cmd.ArgsLenAtDash() == -1 {
				return cmd.Help()
			}

			for _, subCmd := range cmd.Commands() {
				if subCmd.Use == "run" {
					return subCmd.RunE(cmd, args)
				}
			}

			return errors.New("internal failure: failed to find 'run' command")
		},
		SilenceUsage: true,
	}

	rootCmd.PersistentFlags().StringP("dir", "d", ".", "path to dir containing nv.yaml and .env files")
	rootCmd.PersistentFlags().StringP("env", "e", "local", "target environment")

	rootCmd.AddCommand(
		NewPrintCmd(),
		NewResolveCmd(),
		NewRunCmd(),
	)

	return rootCmd
}

func getRunCmd(cmd *cobra.Command) *cobra.Command {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Use == "run" {
			return subCmd
		}
	}

	return nil
}
