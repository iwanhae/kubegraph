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

var _ reconcile.Reconciler = &PodWatcher{}
var _ inject.Client = &PodWatcher{}

type PodWatcher struct {
	Graph *graph.Graph
	client.Reader
}

func (r *PodWatcher) InjectClient(c client.Client) error {
	r.Reader = c
	return nil
}

func (r *PodWatcher) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	id := PodID(req.Namespace, req.Name)

	po := v1.Pod{}
	if err := r.Get(ctx, req.NamespacedName, &po); err != nil {
		if errors.IsNotFound(err) {
			r.Graph.UpdateNode(id, nil, nil)
		}
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	color := "gray"
	if ready := GetPodConditionFromList(po.Status.Conditions, v1.PodScheduled); ready.Status == v1.ConditionTrue {
		color = "yellow"
	}
	if ready := GetPodConditionFromList(po.Status.Conditions, v1.PodInitialized); ready.Status == v1.ConditionTrue {
		color = "skyblue"
	}
	if ready := GetPodConditionFromList(po.Status.Conditions, v1.PodReady); ready.Status == v1.ConditionTrue {
		color = "green"
	}

	edges := []string{
		NodeID(po.Spec.NodeName),
		NSID(req.Namespace),
	}

	for _, ref := range po.ObjectMeta.OwnerReferences {
		if ref.Kind == "ReplicaSet" {
			edges = append(edges, ReplicaSetID(po.Namespace, ref.Name))
		}
	}

	r.Graph.UpdateNode(id, edges,
		Content{
			Color:     color,
			Namespace: po.Namespace,
			Name:      po.Name,
		},
	)
	if ip := po.Status.PodIP; ip != "" {
		ipID := IP(ip)
		edges := []string{id}
		if no := r.Graph.GetNode(ipID); no != nil {
			shouldUpdate := true
			for _, edge := range no.Edges {
				if edge == id {
					shouldUpdate = false
					break
				}
			}
			if shouldUpdate {
				edges = append(no.Edges, id)
			}
		}
		r.Graph.UpdateNode(ipID, edges, Content{Color: "MediumSeaGreen"})
	}
	return reconcile.Result{}, nil
}

func GetPodConditionFromList(conditions []v1.PodCondition, conditionType v1.PodConditionType) v1.PodCondition {
	if conditions == nil {
		return v1.PodCondition{}
	}
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return conditions[i]
		}
	}
	return v1.PodCondition{}
}
