package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(RunCommand(os.Args[1:]))
}

func RunCommand(args []string) int {
	ui := &cli.BasicUi{
		Reader:      bufio.NewReader(os.Stdin),
		Writer:      colorable.NewNonColorable(os.Stdout),
		ErrorWriter: os.Stderr,
	}

	commands := map[string]cli.CommandFactory{
		"get": func() (cli.Command, error) {
			return &GetCommand{
				UI: ui,
			}, nil
		},
		"put": func() (cli.Command, error) {
			return &PutCommand{
				UI: ui,
			}, nil
		},
	}

	cli := &cli.CLI{
		Name:                       "gokv",
		Args:                       args,
		Commands:                   commands,
		HelpFunc:                   cli.BasicHelpFunc("gokv"),
		HelpWriter:                 os.Stderr,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
