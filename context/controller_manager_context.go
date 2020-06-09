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
https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-vsphere/master/pkg/context/controller_manager_context.go

Modifications Copyright 2020 The Cloud Native Scenario Tester Authors.

*/

package context

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/cnstio/controller/record"
)

// ControllerManagerContext is the context of the controller that owns the
// controllers.
type ControllerManagerContext struct {
	context.Context

	// Namespace is the namespace in which the resource is located responsible
	// for running the controller manager.
	Namespace string

	// Name is the name of the controller manager.
	Name string

	// LeaderElectionID is the information used to identify the object
	// responsible for synchronizing leader election.
	LeaderElectionID string

	// LeaderElectionNamespace is the namespace in which the LeaderElection
	// object is located.
	LeaderElectionNamespace string

	// WatchNamespace is the namespace the controllers watch for changes. If
	// no value is specified then all namespaces are watched.
	WatchNamespace string

	// Client is the controller manager's client.
	Client client.Client

	// Logger is the controller manager's logger.
	Logger logr.Logger

	// Recorder is used to record events.
	Recorder record.Recorder

	// Scheme is the controller manager's API scheme.
	Scheme *runtime.Scheme

	// MaxConcurrentReconciles is the maximum number of recocnile requests this
	// controller will receive concurrently.
	MaxConcurrentReconciles int

	genericEventCache sync.Map
}

// String returns ControllerManagerName.
func (c *ControllerManagerContext) String() string {
	return c.Name
}

// GetGenericEventChannelFor returns a generic event channel for a resource
// specified by the provided GroupVersionKind.
func (c *ControllerManagerContext) GetGenericEventChannelFor(gvk schema.GroupVersionKind) chan event.GenericEvent {
	if val, ok := c.genericEventCache.Load(gvk); ok {
		return val.(chan event.GenericEvent)
	}
	val, _ := c.genericEventCache.LoadOrStore(gvk, make(chan event.GenericEvent))
	return val.(chan event.GenericEvent)
}