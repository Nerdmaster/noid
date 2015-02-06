package main

import (
	"fmt"
	"github.com/Nerdmaster/noid/lib/noid"
	"os"
	"strconv"
)

func mintUsageError(message string) {
	fmt.Println(message)
	fmt.Println("")
	mintUsage()
	os.Exit(1)
}

func mintUsage() {
	fmt.Println("Usage: noid-cli mint immediate TEMPLATE SEQUENCE")
	fmt.Println("")
	fmt.Println("Usage: noid-cli mint init TEMPLATE")
	fmt.Println("       noid-cli mint next")
	fmt.Println("")
}

func cmdMintHelp() {
	mintUsage()
	fmt.Println(`Allows minting noids with or without persistence.  If the "immediate"`)
	fmt.Println("sub-command is used, a valid template string and sequence number must also")
	fmt.Println("be present, e.g.:")
	fmt.Println("")
	fmt.Println("    noid-cli mint immediate reedeek 27")
	fmt.Println("")
	fmt.Println(`Using the "next" sub-command requires a noid database in the current working`)
	fmt.Println(`directory.  This can be created with the "init" sub-command, e.g.:`)
	fmt.Println("")
	fmt.Println("    noid-cli mint init reedeek      # Creates noid.db")
	fmt.Println("    noid-cli mint next               # Mints XXXXXX")
	fmt.Println("    noid-cli mint next               # Mints XXXXXX")
	os.Exit(1)
}

// Returns the minted value for the given template at the given sequence.  If
// the template and sequence combination can't be used to generate a minter, an
// error is printed and the application terminates
func mint(template string, sequenceStart uint64) string {
	minter, err := noid.NewSequencedMinter(template, sequenceStart)
	if err != nil {
		mintUsageError(fmt.Sprintf("Error trying to create a minter: %s", err))
	}

	return minter.Mint()
}

func cmdMint(args []string) {
	if len(args) < 1 {
		mintUsageError("Mint command requires a sub-command")
	}

	var fn func([]string)
	argCount := 0

	switch args[0] {
	case "immediate":
		fn = cmdMintImmediate
		argCount = 2

	case "init":
		fn = cmdCreateDatabase
		argCount = 1

	case "next":
		fn = cmdMintNext
	}

	if fn == nil {
		mintUsageError(fmt.Sprintf(`"mint %s" is not a valid command`, args[0]))
	}

	if argCount+1 != len(args) {
		mintUsageError(fmt.Sprintf(`"mint %s" takes %d arguments`, args[0], argCount))
	}

	var newArgs []string
	if len(args) > 1 {
		newArgs = args[1:]
	}
	fn(newArgs)
}

func cmdMintImmediate(args []string) {
	sequenceVal, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		mintUsageError(fmt.Sprintf(`"%s" is not a valid number`, args[1]))
	}
	fmt.Println(mint(args[0], sequenceVal))
}

func cmdCreateDatabase(args []string) {
	fmt.Println("Not implemented")
}

func cmdMintNext([]string) {
	fmt.Println("Not implemented")
}
