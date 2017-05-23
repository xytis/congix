package agent

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mitchellh/cli"
)

type Command struct {
	Ui cli.Ui
}

func (c *Command) Run(args []string) int {
	config := c.readConfig()
	if config == nil {
		return 1
	}


	go c.
	return c.handleSignals(config)
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

func (c *Command) readConfig() *Config {
	flags := flag.NewFlagSet("agent", flag.ContinueOnError)
	flags.Usage = func() { c.Ui.Error(c.Help()) }

	config, err := LoadConfig("./example/config.hcl")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error on parse: %v\n", err))
		return nil
	}
	return config
}

// handleSignals blocks until we get an exit-causing signal
func (c *Command) handleSignals(config *Config) int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGPIPE)

	// Wait for a signal
WAIT:
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	}
	c.Ui.Output(fmt.Sprintf("Caught signal: %v", sig))

	// Skip any SIGPIPE signal (See issue #1798)
	if sig == syscall.SIGPIPE {
		goto WAIT
	}

	// Check if this is a SIGHUP
	/*
		if sig == syscall.SIGHUP {
			if conf := c.handleReload(config); conf != nil {
				*config = *conf
			}
			goto WAIT
		}
	*/

	return 0
}
