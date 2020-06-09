package profiler

import (
	"net/http"
	"net/http/pprof"
)

// RunIfAddressSet runs the profiler if the profiler
// address is provided.
func RunIfAddressSet(o Options) {
	if o.ProfilerAddress != "" {
		o.Logger.Info(
			"Profiler listening for requests",
			"profiler-address", o.ProfilerAddress)
		go runProfiler(o.ProfilerAddress)
	}
}

func runProfiler(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	_ = http.ListenAndServe(addr, mux)
}
