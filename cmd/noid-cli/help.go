package main

import (
	"fmt"
	"os"
)

func helpUsage() {
	fmt.Println("Usage: noid-cli help [command]")
	fmt.Println("")
}

// Shows a list of commands or dispatches to a specific command's help
func cmdHelp(args []string) {
	// Bad syntax
	if len(args) > 1 {
		fmt.Println("Invalid usage - enter a single command to get help")
		fmt.Println("")
		cmdHelpHelp()
	}

	if len(args) == 1 {
		cmd := commands[args[0]]
		if cmd != nil {
			cmd.helpHandler()
			os.Exit(0)
		}

		fmt.Printf("Unknown command, %#v\n", args[0])
		fmt.Println("")
		cmdHelpHelp()
	}

	printUsage()
	for name, cmd := range commands {
		fmt.Printf("%#v: %s\n", name, cmd.helpSummary)
	}

	os.Exit(0)
}

func cmdHelpHelp() {
	helpUsage()
	fmt.Println("Shows help for a specific command, or if no command is present, ")
	fmt.Println("lists all valid commands")
	os.Exit(1)
}
