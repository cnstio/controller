/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Copied from
https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-vsphere/master/pkg/manager/options.go

Modifications Copyright 2020 The Cloud Native Scenario Tester Authors.
*/

package manager

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlmgr "sigs.k8s.io/controller-runtime/pkg/manager"

	// +kubebuilder:scaffold:imports

	"github.com/cnstio/controller/context"
)

// AddToManagerFunc is a function that can be optionally specified with
// the manager's Options in order to explicitly decide what controllers and
// webhooks to add to the manager.
type AddToManagerFunc func(*context.ControllerManagerContext, ctrlmgr.Manager) error

// Options describes the options used to create a new manager.
type Options struct {
	// Prefix for pod, namespace, and id names.
	Prefix string

	// LeaderElectionEnabled is a flag that enables leader election.
	LeaderElectionEnabled bool

	// LeaderElectionID is the name of the config map to use as the
	// locking resource when configuring leader election. Typically,
	// this name is constructed from Prefix and LeaderElectionIDSuffix
	// in defaults() unless defined explicitly.
	LeaderElectionID string

	// LeaderElectionIDSuffix is the suffix name of the config map to
	// use as the locking resource when configuring leader election.
	LeaderElectionIDSuffix string

	// SyncPeriod is the amount of time to wait between syncing the local
	// object cache with the API server.
	SyncPeriod time.Duration

	// MaxConcurrentReconciles the maximum number of allowed, concurrent
	// reconciles.
	//
	// Defaults to the eponymous constant in this package.
	MaxConcurrentReconciles int

	// MetricsAddr is the net.Addr string for the metrics server.
	MetricsAddr string

	// HealthAddr is the net.Addr string for the healthcheck server.
	HealthAddr string

	// LeaderElectionNamespace is the namespace in which the pod running the
	// controller maintains a leader election lock.
	//
	// Default = ""
	LeaderElectionNamespace string

	// PodName is the name of the pod running the controller manager.
	// Typically, this is constructed from the Prefix and PodNameSuffix
	// in the defaults() function unless specified explicitly.
	PodName string

	// PodNameSuffix is the suffix name of the pod running the
	// controller manager.
	//
	// Defaults to the eponymous constant in this package.
	PodNameSuffix string

	// PodNamespace is the namespace in which the pod running the
	// controller is created in. Typically, this is constructed from
	// the Prefix and PodNamespaceSuffix in the defaults() function
	// unless specified explicitly.
	PodNamespace string

	// PodNameSuffix is the namespace suffix in which the pod running the
	// controller is created in.
	//
	// Defaults to the eponymous constant in this package.
	PodNamespaceSuffix string

	// WatchNamespace is the namespace the controllers watch for changes. If
	// no value is specified then all namespaces are watched.
	//
	// Defaults to the eponymous constant in this package.
	WatchNamespace string

	// WebhookPort is the port that the webhook server serves at.
	WebhookPort int

	Logger     logr.Logger
	KubeConfig *rest.Config
	Scheme     *runtime.Scheme
	NewCache   cache.NewCacheFunc

	// AddToManager is a function that can be optionally specified with
	// the manager's Options in order to explicitly decide what controllers
	// and webhooks to add to the manager.
	AddToManager AddToManagerFunc
}

func (o *Options) defaults() {
	if o.Logger == nil {
		o.Logger = ctrllog.Log
	}

	if o.PodName == "" {
		if o.PodNameSuffix == "" {
			if name, err := os.Hostname(); err != nil {
				o.PodName = o.Prefix + name
			} else {
				o.PodName = o.Prefix + "controller-manager"
			}
		} else {
			o.PodName = o.Prefix + o.PodNameSuffix
		}
	}

	if o.SyncPeriod == 0 {
		o.SyncPeriod = DefaultSyncPeriod
	}

	if o.Scheme == nil {
		o.Scheme = runtime.NewScheme()
	}

	if ns, ok := os.LookupEnv("POD_NAMESPACE"); ok {
		o.PodNamespace = ns
	} else if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			o.PodNamespace = ns
		}
	} else {
		if o.PodNamespace == "" {
			if o.PodNamespaceSuffix == "" {
				o.PodNamespace = o.Prefix + "system"
			} else {
				o.PodNamespace = o.Prefix + o.PodNamespaceSuffix
			}
		}
	}

	if o.LeaderElectionID == "" {
		if o.LeaderElectionIDSuffix == "" {
			o.LeaderElectionID = o.Prefix + o.PodName + "-runtime"
		} else {
			o.LeaderElectionID = o.Prefix + o.PodName + o.LeaderElectionIDSuffix
		}
	}

	if o.LeaderElectionNamespace == "" {
		o.LeaderElectionNamespace = "default"
	}
}

// InitFlags initializes the option flags for the manager.
func (o *Options) InitFlags(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}

	flag.StringVar(
		&o.MetricsAddr,
		"metrics-addr",
		":8080",
		"The address the metric endpoint binds to.")
	flag.BoolVar(
		&o.LeaderElectionEnabled,
		"enable-leader-election",
		true,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(
		&o.LeaderElectionID,
		"leader-election-id",
		"",
		"Name of the config map to use as the locking resource when configuring leader election.")
	flag.StringVar(
		&o.LeaderElectionNamespace,
		"leader-election-namespace",
		"",
		"Name of the namespace to use for the configmap locking resource when configuring leader election.")
	flag.StringVar(
		&o.WatchNamespace,
		"namespace",
		"",
		"Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.")
	flag.DurationVar(
		&o.SyncPeriod,
		"sync-period",
		DefaultSyncPeriod,
		"The interval at which cluster-api objects are synchronized")
	flag.IntVar(
		&o.MaxConcurrentReconciles,
		"max-concurrent-reconciles",
		10,
		"The maximum number of allowed, concurrent reconciles.")
	flag.StringVar(
		&o.PodNameSuffix,
		"pod-name-suffix",
		"controller-manager",
		"The suffix name of the pod running the controller manager.")
	flag.StringVar(
		&o.PodNamespaceSuffix,
		"pod-namespace-suffix",
		"controller-manager",
		"The suffix name of the pod namespace running the controller manager.")
	flag.IntVar(
		&o.WebhookPort,
		"webhook-port",
		DefaultWebhookServiceContainerPort,
		"Webhook Server port (set to 0 to disable)")
	flag.StringVar(
		&o.HealthAddr,
		"health-addr",
		":9440",
		"The address the health endpoint binds to.",
	)
}
