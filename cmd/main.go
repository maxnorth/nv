package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/maxnorth/nv/resolver"
)

func main() {
	var optionArgs []string
	var commandArgs []string
	var command string
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if command != "" {
			commandArgs = append(commandArgs, arg)
		} else if arg == "--" || !strings.HasPrefix(arg, "-") {
			command = arg
		} else {
			optionArgs = append(optionArgs, arg)
		}
	}

	r := resolver.Load()
	values := r.Resolve()

	if command == "" {
		printEnv(values, optionArgs)
		return
	}
	if command == "shell" {
		command = "--"
		commandArgs = append([]string{os.Getenv("SHELL")}, commandArgs...)
	}
	if command == "--" {
		execCommand(commandArgs)
		return
	}

	fmt.Println("unrecognized command")
	os.Exit(1)
}

func execCommand(commandArgs []string) {
	if len(commandArgs) == 0 {
		fmt.Println("command was not provided")
		os.Exit(1)
	}

	command, _ := commandArgs[0], []string{}
	// if len(commandArgs) > 1 {
	// 	args = commandArgs[1:]
	// }
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

func printEnv(values map[string]string, optionArgs []string) {
	output := "dotenv"
	for i, arg := range optionArgs {
		if arg == "-o" || arg == "--output" {
			if len(optionArgs) > i+1 {
				output = optionArgs[i+1]
				break
			} else {
				// TODO
				fmt.Println("missing value for --output arg")
				os.Exit(1)
			}
		}
	}

	var outputTemplate string
	switch output {
	case "yaml":
		outputTemplate = "%s: \"%s\"\n"
	case "dotenv":
		outputTemplate = "%s=%s\n"
	case "docker":
		outputTemplate = "--env %s=\"%s\"\n"
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

// replace env vars starting with @
//  - resolve the provider
//    - if provider not found, fail w/ helpful message
//    - later provide an option to ignore specific values or any not found value
//  - if no : the key is the map target
//    - if : use that as the key
//  - using the key, invoke the provider and resolve the value
//  - map the value to the target env var

// need to identify resolution order and ordering issues
