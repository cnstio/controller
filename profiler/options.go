package profiler

import (
	"flag"
	"os"

	"github.com/go-logr/logr"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// Options for configuring the http profiler.
type Options struct {
	ProfilerAddress string
	Logger          logr.Logger
}

func (o *Options) defaults() {
	if o.Logger == nil {
		o.Logger = ctrllog.Log.WithName("profiler")
	}
}

// InitFlags initializes the profiler flags.
func (o *Options) InitFlags(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}

	fs.StringVar(
		&o.ProfilerAddress,
		"profiler-address",
		os.Getenv("PROFILER_ADDR"),
		"Bind address to expose the pprof profiler (e.g. localhost:6060)")
}
