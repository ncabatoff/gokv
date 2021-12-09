package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/philippgille/gokv/encoding"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mitchellh/cli"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/bbolt"
	"github.com/posener/complete"
)

var _ cli.Command = (*GetCommand)(nil)

type GetCommand struct {
	UI         *cli.BasicUi
	flagBucket string
	flagCodec  string
	flagDriver string

	doneCtx   context.Context
	cancelCtx context.CancelFunc
	wg        sync.WaitGroup
}

func (c *GetCommand) Synopsis() string {
	return "Fetch a value from a KV store"
}

func (c *GetCommand) Help() string {
	helpText := `
Usage: gokv get [-driver=bbolt] <source> <key>
`
	return strings.TrimSpace(helpText)
}

func (c *GetCommand) Flags() *flag.FlagSet {
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

func (c *GetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *GetCommand) AutocompleteFlags() complete.Flags {
	return nil
	//return c.Flags().Completions()
}

func (c *GetCommand) Run(args []string) int {
	c.doneCtx, c.cancelCtx = context.WithCancel(context.Background())
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
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", len(args)))
		return 1
	}

	location, key := args[0], args[1]
	store, err := getStore(c.flagDriver, location, c.flagBucket)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error opening %s store: %v", c.flagDriver, err))
		return 1
	}

	v, err := getValue(store, key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading key %q: %v", key, err))
		return 1
	}

	// Don't use println because we don't want to be appending newlines to
	// binary data
	fmt.Print(string(v))

	defer func() {
		c.cancelCtx()
		c.wg.Wait()
	}()

	return 0
}

func getValue(store gokv.Store, key string) ([]byte, error) {
	var b []byte
	found, err := store.Get(key, &b)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("key %q doesn't exist", key)
	}
	return b, nil
}

func getStore(driver, location, bucket string) (gokv.Store, error) {
	switch driver {
	case "bbolt":
		options := bbolt.DefaultOptions
		options.Path = location
		options.BucketName = bucket
		options.Codec = encoding.NoCodec{}
		return bbolt.NewStore(options)
	default:
		return nil, fmt.Errorf("unknown driver %q", driver)
	}
}
