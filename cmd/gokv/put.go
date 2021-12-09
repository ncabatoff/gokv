package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mitchellh/cli"
	"github.com/philippgille/gokv"
	"github.com/posener/complete"
)

var _ cli.Command = (*PutCommand)(nil)

type PutCommand struct {
	UI         *cli.BasicUi
	flagBucket string
	flagCodec  string
	flagDriver string

	doneCtx   context.Context
	cancelCtx context.CancelFunc
	wg        sync.WaitGroup
}

func (c *PutCommand) Synopsis() string {
	return "Store a value into a KV store"
}

func (c *PutCommand) Help() string {
	helpText := `
Usage: gokv put [-driver=bbolt] <source> <key> <value>
`
	return strings.TrimSpace(helpText)
}

func (c *PutCommand) Flags() *flag.FlagSet {
	set := flag.NewFlagSet("flags", flag.ContinueOnError)

	set.StringVar(
		&c.flagDriver,
		"driver",
		"bbolt",
		"Store type, one of 'bbolt'",
	)
	set.StringVar(
		&c.flagBucket,
		"bucket",
		"bucket",
		"",
	)
	set.StringVar(
		&c.flagCodec,
		"codec",
		"none",
		"Byte encoding for values, one of 'none'",
	)
	return set
}

func (c *PutCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *PutCommand) AutocompleteFlags() complete.Flags {
	return nil
	//return c.Flags().Completions()
}

func (c *PutCommand) Run(args []string) int {
	// handle ctrl-c while waiting for the callback
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	go func() {
		<-sigCh
		c.cancelCtx()
	}()

	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2-3, got %d)", len(args)))
		return 1
	case len(args) > 3:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2-3, got %d)", len(args)))
		return 1
	}

	location, key := args[0], args[1]
	var value []byte
	if len(args) > 2 {
		value = []byte(args[2])
	} else {
		b, err := ioutil.ReadAll(c.UI.Reader)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error reading value from stdin: %v", err))
			return 1
		}
		value = b
	}

	store, err := getStore(c.flagDriver, location, c.flagBucket)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error opening %s store: %v", c.flagDriver, err))
		return 1
	}

	err = putValue(store, key, value)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing key %q: %v", key, err))
		return 1
	}

	return 0
}

func putValue(store gokv.Store, key string, value []byte) error {
	return store.Set(key, value)
}
