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

type NamespaceWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *NamespaceWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *NamespaceWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := NSID(req.Name)
	ns := v1.Namespace{}
	if err := r.Get(ctx, req.NamespacedName, &ns); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	edges := []string{}

	r.Graph.UpdateNode(id, edges, Content{
		Color:     "black",
		Name:      ns.Name,
		Namespace: ns.Namespace,
	})
	return reconcile.Result{}, nil
}
