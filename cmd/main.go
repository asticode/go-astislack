package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astislack"
	"github.com/asticode/go-astitools/config"
	"github.com/asticode/go-astitools/http"
	"github.com/pkg/errors"
)

var (
	flagConfig = flag.String("c", "", "the config path")
	flagUser   = flag.String("u", "", "the user")
)

func main() {
	// Parse flags
	flag.Parse()
	astilog.FlagInit()

	// Create configuration
	c, err := newConfiguration()
	if err != nil {
		astilog.Fatalf("main: creating configuration failed")
	}

	// Create logger
	astilog.SetLogger(astilog.New(c.Logger))

	// Create slack
	s := astislack.New(c.Slack)

	// Get user
	user := *flagUser
	if user == "" {
		// Get me
		var me astislack.Me
		if me, err = s.Me(); err != nil {
			astilog.Fatal(errors.Wrap(err, "main: getting me failed"))
		}
		user = me.User
	}

	// No user
	if user == "" {
		astilog.Fatal("main: no user found")
	}
	astilog.Infof("main: processing user '%s'", user)

	// Delete
	if err = s.Delete(fmt.Sprintf("from:%s before:today", user)); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: deleting failed"))
	}
}

type Configuration struct {
	Logger astilog.Configuration   `toml:"logger"`
	Slack  astislack.Configuration `toml:"slack"`
}

func newConfiguration() (*Configuration, error) {
	i, err := asticonfig.New(&Configuration{
		Logger: astilog.Configuration{
			AppName: "astislack",
			Format:  astilog.FormatText,
			Out:     astilog.OutStdOut,
		},
		Slack: astislack.Configuration{Sender: astihttp.SenderOptions{
			RetrySleep: 5 * time.Second,
			RetryMax:   5,
		}},
	}, *flagConfig, &Configuration{
		Logger: astilog.FlagConfig(),
		Slack:  astislack.FlagConfig(),
	})
	return i.(*Configuration), err
}
