package cmd

import (
	"fmt"

	"github.com/maxnorth/nv/resolver"
	"github.com/spf13/cobra"
)

func ResolveCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use: "resolve",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := resolver.Load(cmd.Flag("env").Value.String())

			value, err := r.ResolveUrl(args[0])
			if err != nil {
				return err
			}

			fmt.Println(value)

			return nil
		},
	}

	return runCmd
}
