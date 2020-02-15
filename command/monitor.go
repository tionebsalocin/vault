package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// ==================
	var logCh chan string
	stopCh := make(chan struct{})
	defer close(stopCh)

	// TODO: fix this so I can pass query options for log level
	//logCh, err = client.Sys().Monitor("INFO", false, eventDoneCh, nil)
	logCh, err = client.Sys().Monitor("INFO", false, stopCh, false)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error starting monitor: %s", err))
		return 1
	}

	fmt.Println("creating the signal channel")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	z, _ := os.Create("good-times.txt")

	go func() {
		//defer close(eventDoneCh)
	OUTER:
		for {
			select {
			case log := <-logCh:
				z.WriteString(fmt.Sprintf("got a log: %v\n", log))
				if log == "" {
					z.WriteString("breaking cuz ya\n")
					break OUTER
				}
				c.UI.Info(log)
			}
		}

		//c.lock.Lock()
		//defer c.lock.Unlock()
		//if !c.quitting {
		//	c.UI.Info("")
		//	c.UI.Output("Remote side ended the monitor! This usually means that the\n" +
		//		"remote side has exited or crashed.")
		//}
	}()

	fmt.Println("about to wait on shit")
	select {
	//case <-eventDoneCh:
	//	fmt.Println("got something on eventdonech")
	//	return 1
	case <-signalCh:
		fmt.Println("interrupt received")
		return 0
		//case <-c.shutdownCh:
		//	c.lock.Lock()
		//	c.quitting = true
		//	c.lock.Unlock()
	}

	fmt.Println("how did we get here?  isn't select supposed to block?")

	return 0
	// ==========================

	//stopCh := make(chan struct{})
	//defer close(stopCh)
	//
	//respCh, errCh := client.Sys().Monitor(stopCh)
	//
	//select {
	//case err := <-errCh:
	//	c.UI.Error(fmt.Sprintf("Error starting monitor: %s", err))
	//	return 1
	//default:
	//}
	//
	//signalCh := make(chan os.Signal, 1)
	//signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	//
	//i := 0
	//
	////for {
	////	select {
	////	case log := <-frames:
	////		if log == nil {
	////			continue
	////		}
	////		if strings.Contains(string(log.Data), "[DEBUG]") {
	////			break OUTER
	////		}
	////	case err := <-errCh:
	////		t.Fatalf("error: %v", err)
	////	case <-time.After(2 * time.Second):
	////		require.Fail(t, "failed to get a DEBUG log message")
	////	}
	////}
	//
	////go func() {
	//	for {
	//		fmt.Printf("i = %v\n", i)
	//		select {
	//		case mResp := <-respCh:
	//			if mResp == nil {
	//				continue
	//			} else {
	//				fmt.Println(mResp)
	//			}
	//		case <-signalCh:
	//			fmt.Println("interrupt received")
	//			return 0
	//		}
	//		i++
	//
	//		if i > 10 {
	//			break
	//		}
	//	}
	////}()
	//
	//return 0
}
