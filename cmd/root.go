package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "nv",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			dir := cmd.Flag("dir").Value.String()
			workDir, _ := os.Getwd()
			dirPath := path.Join(workDir, dir)
			if err := os.Chdir(dirPath); err != nil {
				// TODO: other kinds of errors? permission issues?
				return fmt.Errorf("target directory not found: %s", dirPath)
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

			return errors.New("internal error: failed to find 'run' command")
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
