package cmd

import (
	"fmt"
	"os"

	"github.com/maxnorth/nv/resolver"
	"github.com/spf13/cobra"
)

func PrintCmd() *cobra.Command {
	printCmd := &cobra.Command{
		Use: "print",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := resolver.Load(cmd.Flag("env").Value.String())

			values, err := r.ResolveEnvironment()
			if err != nil {
				return err
			}

			printEnv(values, cmd.Flag("output").Value.String())

			return nil
		},
	}
	printCmd.Flags().StringP("output", "o", "env", "output format, supports 'env', 'json', or 'yaml'")

	return printCmd
}

func printEnv(values map[string]string, output string) {
	if output == "" {
		fmt.Println("missing value for --output arg")
		os.Exit(1)
	}

	var outputTemplate string
	switch output {
	case "yaml":
		outputTemplate = "%s: \"%s\"\n"
	case "env":
		outputTemplate = "%s=%s\n"
	case "shell":
		outputTemplate = "export %s=\"%s\"\n"
	}

	if outputTemplate != "" {
		for key, value := range values {
			fmt.Printf(outputTemplate, key, value)
		}
		return
	}

	switch output {
	case "keys":
		outputTemplate = "%s\n"
	}

	if outputTemplate != "" {
		for key, _ := range values {
			fmt.Printf(outputTemplate, key)
		}
		return
	}
}
