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

type ServiceWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *ServiceWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *ServiceWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := SvcID(req.Namespace, req.Name)
	svc := v1.Service{}
	if err := r.Get(ctx, req.NamespacedName, &svc); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	edges := []string{
		EPID(svc.Namespace, svc.Name),
		NSID(req.Namespace),
	}

	r.Graph.UpdateNode(id, edges, Content{
		Color: "orange",
	})

	r.Graph.UpdateNode(IP(svc.Spec.ClusterIP), []string{id}, Content{
		Color: "SlateBlue",
	})
	return reconcile.Result{}, nil
}
