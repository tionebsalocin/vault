package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*MonitorCommand)(nil)
var _ cli.CommandAutocomplete = (*MonitorCommand)(nil)

type MonitorCommand struct {
	*BaseCommand

	flagLogLevel string
}

func (c *MonitorCommand) Synopsis() string {
	return "Stream log messages from a Vault server"
}

func (c *MonitorCommand) Help() string {
	helpText := `
Usage: vault monitor [options]

	Stream log messages of a Vault server. The monitor command lets you listen
	for log levels that may be filtered out of the server logs. For example,
	the server may be logging at the INFO level, but with the monitor command
	you can set -log-level DEBUG.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *MonitorCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetNone)

	f := set.NewFlagSet("Monitor Options")
	f.StringVar(&StringVar{
		Name:       "log-level",
		Target:     &c.flagLogLevel,
		Default:    "info",
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "error"),
		Usage: "If passed, the log level to monitor logs. Supported values" +
			"(in order of detail) are \"trace\", \"debug\", \"info\", \"warn\"" +
			" and \"error\".",
	})

	return set
}

func (c *MonitorCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *MonitorCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *MonitorCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	fmt.Println("args")
	fmt.Println(args)

	//client, err := c.Client()
	//if err != nil {
	//	c.UI.Error(err.Error())
	//	return 2
	//}

	return 0
}
