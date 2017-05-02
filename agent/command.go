package agent

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/xytis/congix/structs/config"
)

type Command struct {
	Ui cli.Ui
}

func (c *Command) Run(args []string) int {

	// Make a new, empty config.
	config := &Config{
		Nginx:  &config.NginxConfig{},
		Consul: &config.ConsulConfig{},
	}
	fmt.Printf("config: %v\n", config)

	flags := flag.NewFlagSet("agent", flag.ContinueOnError)
	flags.Usage = func() { c.Ui.Error(c.Help()) }

	fmt.Println("parsing file")
	_, err := ParseFile("./example/config.hcl")
	fmt.Errorf("error on parse: %v\n", err)

	fmt.Println("asking for help")
	return cli.RunResultHelp
}

func (c *Command) Synopsis() string {
	return "Runs a Congix agent"
}

func (c *Command) Help() string {
	helpText := `
Usage: congix agent [options]

  Starts the Congix agent and runs until an interupt is received.

`

	return strings.TrimSpace(helpText)
}
