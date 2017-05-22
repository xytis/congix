package main

import (
	//"flag"
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/xytis/congix/agent"
)

var (
	VERSION = "0.0.1"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {
	c := cli.NewCLI("congix", VERSION)
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	c.Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &agent.Command{
				Ui: ui,
			}, nil
		},
		"check": func() (cli.Command, error) {
			return &agent.Command{}, nil
		},
	}
	c.Args = args

	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
