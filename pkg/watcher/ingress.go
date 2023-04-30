package watcher

import (
	"context"

	"github.com/iwanhae/kubegraph/pkg/graph"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

var _ reconcile.Reconciler = &NodeWatcher{}
var _ inject.Client = &NodeWatcher{}

type IngressWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *IngressWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *IngressWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := IngressID(req.Namespace, req.Name)
	ing := networkingv1.Ingress{}
	if err := r.Get(ctx, req.NamespacedName, &ing); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	edges := []string{
		NSID(req.Namespace),
	}

	for _, rule := range ing.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			edges = append(edges, SvcID(ing.Namespace, path.Backend.Service.Name))
		}
	}

	r.Graph.UpdateNode(id, edges, Content{
		Color: "violet",
		Name:  ing.Name,
	})
	return reconcile.Result{}, nil
}
