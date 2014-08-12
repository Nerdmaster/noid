package main

type Command struct {
	handler func([]string)
	helpHandler func()
	helpSummary string
}

var commands = make(map[string]*Command)

func initCommands() {
	commands["help"] = &Command{handler: cmdHelp, helpHandler: cmdHelpHelp, helpSummary: "Displays this usage page"}
}
