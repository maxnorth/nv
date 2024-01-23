package main

import (
	"github.com/maxnorth/nv/cmd"
)

func main() {
	cmd.RootCmd().Execute()
}

// func main2() {
// 	var optionArgs []string
// 	var commandArgs []string
// 	var command string
// 	for i, arg := range os.Args {
// 		if i == 0 {
// 			continue
// 		}
// 		if command != "" {
// 			commandArgs = append(commandArgs, arg)
// 		} else if arg == "--" || !strings.HasPrefix(arg, "-") {
// 			command = arg
// 		} else {
// 			optionArgs = append(optionArgs, arg)
// 		}
// 	}

// 	r := resolver.Load()
// 	values := r.Resolve()

// 	if command == "" {
// 		printEnv(values, "")
// 		return
// 	}
// 	if command == "shell" {
// 		command = "--"
// 		commandArgs = append([]string{os.Getenv("SHELL")}, commandArgs...)
// 	}
// 	if command == "--" {
// 		execCommand(commandArgs)
// 		return
// 	}

// 	fmt.Println("unrecognized command")
// 	os.Exit(1)
// }
