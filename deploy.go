package apps

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/hooklift/lift/ui"
	"github.com/hooklift/sync"
	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
)

type cmdOpts struct {
	appName    string
	appDir     string
	syncClient sync.SyncClient
	logOutput  io.Writer
}

// CmdOption defines a function type for passing command options.
type CmdOption func(*cmdOpts)

// WithApp allows to explicitly specified the app to deploy.
func WithApp(name string) CmdOption {
	return func(o *cmdOpts) {
		o.appName = name
	}
}

// WithAppDir specifies the root of the app's sources to sync to the server.
func WithAppDir(dpath string) CmdOption {
	return func(o *cmdOpts) {
		o.appDir = dpath
	}
}

// WithSyncClient sets the gRPC client to use in order to sync app sources to the server and deploy the app.
func WithSyncClient(client sync.SyncClient) CmdOption {
	return func(o *cmdOpts) {
		o.syncClient = client
	}
}

// WithLogOutput allows to set a writer where all the logging will be done. By default it uses os.Stdout.
func WithLogOutput(w io.Writer) CmdOption {
	return func(o *cmdOpts) {
		o.logOutput = w
	}
}

// Deploy syncs app source, builds, releases and deploys an app onto Hooklift.
func Deploy(ctx context.Context, opts ...CmdOption) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrapf(err, "failed getting current working dir")
	}

	options := &cmdOpts{
		logOutput: os.Stdout,
		appDir:    cwd,
	}
	for _, o := range opts {
		o(options)
	}

	if options.syncClient == nil {
		return errors.New("sync client required")
	}

	// Read metadata to grab app ID, just in case the app has been previously deployed and configured.
	md := new(DeployMetadata)
	if err := md.Read(); err != nil && options.appName == "" {
		// TODO(c4milo): Use logWriter instead?
		ui.Debug("App metadata doesn't seem to exist, we will create a new app: %s", err.Error())
	}

	syncID := uuid.NewV4().String()
	syncOpts := []sync.Option{
		sync.WithCacheID(md.CacheID),
		sync.WithSyncID(syncID),
		sync.WithClient(options.syncClient),
		sync.WithProgressCb(func(sent uint) {
			// Update global counter so global progress bar moves.
		}),
	}

	ui.Info("Syncing source files... ")

	// TODO(c4milo): Create buffered channel with 5 elements of capacity in order
	// to upload 5 files at a time.
	// Walk app folder structure
	filepath.Walk(options.appDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed walking app source")
		}

		// Ignore directories and let the server re-create them using the file's full path.
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "failed opening %q", path)
		}

		defer func() {
			if err := f.Close(); err != nil {
				ui.Error("failed closing %q: %v", path, err)
			}
		}()

		// Each sync call is made over a different blocking HTTP2 stream.
		// TODO(c4milo): Parallelize this to upload 5 files at the same time. (use sync/errgroup)
		return sync.Push(ctx, f, syncOpts...)
	})

	md.CacheID = syncID
	if err := md.Write(); err != nil {
		return errors.Wrapf(err, "failed writing deployment medatata")
	}
	ui.Info("done\n")

	res, err := options.syncClient.Deploy(ctx, &sync.DeployRequest{
		AppId:  md.AppID,
		SyncId: syncID,
	})

	if err != nil {
		return errors.Wrapf(err, "failed deploying")
	}

	_ = res.OutputUrl
	// TODO(c4milo): gRPC request to output URL in order to read deploy logs?
	// TODO(c4milo): From server, close connection once deployment is done.

	return nil
}
