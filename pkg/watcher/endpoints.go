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

type EPWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *EPWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *EPWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := EPID(req.Namespace, req.Name)
	ep := v1.Endpoints{}
	if err := r.Get(ctx, req.NamespacedName, &ep); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	edges := []string{
		NSID(req.Namespace),
	}

	for _, subset := range ep.Subsets {
		for _, address := range subset.Addresses {
			edges = append(edges, IP(address.IP))
		}
	}
	edges = append(edges, SvcID(ep.Namespace, ep.Name))

	r.Graph.UpdateNode(id, edges, Content{
		Color: "tomato",
		Name:  ep.Name,
	})
	return reconcile.Result{}, nil
}
