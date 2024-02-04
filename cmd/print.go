package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/maxnorth/nv/internal/resolver"
	"github.com/spf13/cobra"
)

func NewPrintCmd() *cobra.Command {
	printCmd := &cobra.Command{
		Use: "print",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := resolver.Load(cmd.Flag("env").Value.String())
			if err != nil {
				return err
			}

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

func printEnv(values map[string]string, output string) error {
	if output == "" {
		return errors.New("missing value for --output arg")
	}

	if output == "json" {
		jsonOutput, err := json.MarshalIndent(values, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonOutput))
		return nil
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
		return nil
	}

	switch output {
	case "keys":
		outputTemplate = "%s\n"
	}

	if outputTemplate != "" {
		for key, _ := range values {
			fmt.Printf(outputTemplate, key)
		}
		return nil
	}

	return nil
}
