package main

import (
	"context"
	"os"

	docopt "github.com/docopt/docopt-go"
	"github.com/hooklift/lift/config"
	"github.com/hooklift/lift/ui"
	"github.com/hooklift/sync"
	"github.com/lift-plugins/apps"
	"github.com/lift-plugins/auth/openidc/grpcutil"
)

var (
	// Version is defined in compilation time.
	Version string

	// Name is defined in compilation time.
	Name string
)

const usage = `
Deploy and manage your apps in Hooklift.

Usage:
  apps deploy [--app=APPNAME]
  apps destroy [--app=APPNAME]
  apps logs [--app=APPNAME] [--tail] [--num=LINES]
  apps info [--app=APPNAME]
  apps rename <newname> [--app=APPNAME]
  apps transfer <recipient> [--app=APPNAME]
  apps open [--app=APPNAME]
  apps errors [--app=APPNAME]
  apps -h | --help
  apps -v | --verbose
  apps --version

Commands:
  deploy                                   Deploys an app.
  destroy                                  Removes the app, it cannot be recovered.
  logs                                     Prints out app logs.
  info                                     Shows app information.
  rename                                   Renames the app.
  transfer                                 Transfer app ownership.
  open                                     Opens app in a web browser.
  errors                                   Shows app errors.

Options:
  -a --app=APPNAME                         Application name to use when running a command.
  -t --tail                                Tails app logs in soft realtime.
  -n --num=LINES                           Number of lines to tail from logs [default: 200].
  -h --help                                Shows this screen.
  -v --version                             Shows version of this plugin.
`

func main() {
	args, err := docopt.Parse(usage, nil, false, "", false, false)
	if err != nil {
		ui.Debug("docopt failed to parse command: ->%#v<-", err)
		ui.Info(usage)
		os.Exit(1)
	}

	ctx := context.Background()

	if args["--version"].(bool) {
		ui.Info(Version)
		return
	}

	if args["--help"].(bool) {
		ui.Info(usage)
		return
	}

	if args["deploy"].(bool) {
		deploy(ctx, args)
		return
	}

	if args["destroy"].(bool) {
		destroy(ctx, args)
		return
	}

	if args["logs"].(bool) {
		logs(ctx, args)
		return
	}

	if args["info"].(bool) {
		info(ctx, args)
		return
	}

	if args["rename"].(bool) {
		rename(ctx, args)
		return
	}

	if args["transfer"].(bool) {
		transfer(ctx, args)
		return
	}

	if args["open"].(bool) {
		open(ctx, args)
		return
	}

	if args["errors"].(bool) {
		errors(ctx, args)
		return
	}
}

func deploy(ctx context.Context, args map[string]interface{}) {
	var appName string

	val, ok := args["--app"]
	if ok && val != nil {
		appName = val.(string)
	}

	// TODO(c4milo): use plugin version as part of the user-agent too.
	syncConn, err := grpcutil.Connection(config.SyncURI, Name)
	if err != nil {
		ui.Debug("%+v", err)
		ui.Error("failed connecting to Hooklift Sync service: %s", err.Error())
		return
	}

	opts := []apps.CmdOption{
		apps.WithApp(appName),
		apps.WithSyncClient(sync.NewSyncClient(syncConn)),
	}
	if err := apps.Deploy(ctx, opts...); err != nil {
		ui.Debug("%+v", err)
		ui.Error("failed deploying app: %s", err.Error())
	}
}

func destroy(ctx context.Context, args map[string]interface{}) {
	// app := args["--app"].(string)
}

func logs(ctx context.Context, args map[string]interface{}) {

}

func info(ctx context.Context, args map[string]interface{}) {
	// app := args["--app"].(string)
}

func rename(ctx context.Context, args map[string]interface{}) {
	// app := args["--app"].(string)
}

func transfer(ctx context.Context, args map[string]interface{}) {
	// app := args["--app"].(string)
}

func open(ctx context.Context, args map[string]interface{}) {
	// app := args["--app"].(string)
}

func errors(ctx context.Context, args map[string]interface{}) {
	// app := args["--app"].(string)
}
