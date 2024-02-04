package cmd

import (
	"fmt"

	"github.com/maxnorth/nv/internal/resolver"
	"github.com/spf13/cobra"
)

func NewResolveCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use: "resolve",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := resolver.Load(cmd.Flag("env").Value.String())
			if err != nil {
				return err
			}

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
