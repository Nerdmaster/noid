package main

import (
	"fmt"
	"nerdbucket.com/go/noid/noid"
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
	fmt.Println("    noid-cli mint init reedeek       # Creates noid.db with serialized minter")
	fmt.Println("    noid-cli mint next               # Prints out q67j4g")
	fmt.Println("    noid-cli mint next               # Prints out y67j4r")
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

func mintFromArgs(args []string) string {
	sequenceVal, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		mintUsageError(fmt.Sprintf(`Unable to mint: sequence value "%s" is not a number`, args[1]))
	}
	return mint(args[0], sequenceVal)
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
	fmt.Println(mintFromArgs(args))
}

func cmdCreateDatabase(args []string) {
	// First make sure we can create the file - it MUST be a new file
	f, err := os.OpenFile("noid.db", os.O_CREATE|os.O_EXCL|os.O_RDWR, 0660)
	if err != nil {
		mintUsageError(fmt.Sprintf("Unable to create noid.db: %s", err))
	}
	defer f.Close()

	// Now make sure the template is legit
	template := args[0]
	m, err := noid.NewMinter(template)
	if err != nil {
		mintUsageError(fmt.Sprintf("Unable to create a minter with template %s: %s", err))
	}

	m.WriteJSON(f)
}

func cmdMintNext([]string) {
	m, err := noid.NewMinterFromJSONFile("noid.db")
	if err != nil {
		mintUsageError(fmt.Sprintf("Error building minter from noid.db: %s", err))
	}
	fmt.Println(m.Mint())

	f, err := os.OpenFile("noid.db", os.O_RDWR, 0660)
	if err != nil {
		mintUsageError(fmt.Sprintf("Unable to re-serialize minter: %s", err))
	}
	defer f.Close()
	m.WriteJSON(f)
}
