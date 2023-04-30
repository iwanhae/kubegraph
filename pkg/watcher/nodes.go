package watcher

import (
	"context"

	"github.com/iwanhae/kubegraph/pkg/graph"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

var _ reconcile.Reconciler = &NodeWatcher{}
var _ inject.Client = &NodeWatcher{}

type NodeWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *NodeWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *NodeWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := NodeID(req.Name)
	no := v1.Node{}
	if err := r.Get(ctx, req.NamespacedName, &no); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	edges := []string{}
	for _, addr := range no.Status.Addresses {
		edges = append(edges, IP(addr.Address))
	}

	r.Graph.UpdateNode(id, edges, Content{
		Color: "cyan",
		Name:  no.Name,
	})

	return reconcile.Result{}, nil
}
