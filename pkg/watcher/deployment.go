package watcher

import (
	"context"

	"github.com/iwanhae/kubegraph/pkg/graph"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

var _ reconcile.Reconciler = &NodeWatcher{}
var _ inject.Client = &NodeWatcher{}

type DeploymentWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *DeploymentWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *DeploymentWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := DeploymentID(req.Namespace, req.Name)
	deploy := appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, &deploy); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	edges := []string{
		NSID(req.Namespace),
	}

	r.Graph.UpdateNode(id, edges, Content{
		Color:     "darkblue",
		Name:      deploy.Name,
		Namespace: deploy.Namespace,
	})
	return reconcile.Result{}, nil
}
