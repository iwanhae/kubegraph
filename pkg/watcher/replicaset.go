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

type ReplicasetWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *ReplicasetWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *ReplicasetWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := ReplicaSetID(req.Namespace, req.Name)
	rs := appsv1.ReplicaSet{}
	if err := r.Get(ctx, req.NamespacedName, &rs); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	edges := []string{NSID(req.Namespace)}
	for _, ref := range rs.ObjectMeta.OwnerReferences {
		if ref.Kind == "Deployment" {
			edges = append(edges, DeploymentID(rs.Namespace, ref.Name))
		}
	}

	r.Graph.UpdateNode(id, edges, Content{
		Color:     "darkcyan",
		Name:      rs.Name,
		Namespace: rs.Namespace,
	})
	return reconcile.Result{}, nil
}
