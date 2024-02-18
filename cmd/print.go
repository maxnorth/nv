package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/alessio/shellescape.v1"
	"gopkg.in/yaml.v3"

	"github.com/maxnorth/nv/internal/resolver"
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

			keys := r.LoadedVars
			for _, key := range r.LoadedVars {
				delete(values, key)
			}
			for key := range values {
				keys = append(keys, key)
			}

			sort.Strings(keys)

			printEnv(keys, cmd.Flag("output").Value.String())

			return nil
		},
	}
	printCmd.Flags().StringP("output", "o", "env", "output format, supports 'env', 'json', or 'yaml'")

	return printCmd
}

func printEnv(keys []string, output string) error {
	if output == "" {
		return errors.New("missing value for --output arg")
	}

	values := map[string]string{}
	for _, key := range keys {
		values[key] = os.Getenv(key)
	}

	if output == "json" {
		jsonOutput, err := json.MarshalIndent(values, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonOutput))
		return nil
	}

	if output == "yaml" {
		yamlOutput, err := yaml.Marshal(values)
		if err != nil {
			return err
		}
		fmt.Print(string(yamlOutput))
		return nil
	}

	if output == "shell" {
		for _, key := range keys {
			fmt.Printf("export %s=%s\n", key, shellescape.Quote(values[key]))
		}
		return nil
	}

	if output == "env" {
		for _, key := range keys {
			fmt.Printf("%s=%s\n", key, values[key])
		}
		return nil
	}

	if output == "keys" {
		for _, key := range keys {
			fmt.Println(key)
		}
		return nil
	}

	return nil
}
