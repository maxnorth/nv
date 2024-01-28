package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "nv",
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
	}
	rootCmd.PersistentFlags().StringP("env", "e", "local", "target environment")

	rootCmd.AddCommand(
		PrintCmd(),
		ResolveCmd(),
		RunCmd(),
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
