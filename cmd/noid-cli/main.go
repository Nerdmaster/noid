package main

import (
	"os"
	"fmt"
)

func printUsage() {
	fmt.Println("Usage: noid-cli [command] ...")
	fmt.Println("")
}

func printUsageOnError() {
	printUsage()
	fmt.Println("Type `noid-cli help` for a list of valid commands")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsageOnError()
	}

	command := os.Args[1]
	other := os.Args[2:]

	initCommands()
	cmd := commands[command]
	if cmd != nil {
		cmd.handler(other)
		os.Exit(0)
	}

	fmt.Printf("Invalid command %#v\n", command)
	printUsageOnError()
}
