package main

import (
	"context"
	"os"
	"time"

	"github.com/iwanhae/kubegraph/pkg/graph"
	"github.com/iwanhae/kubegraph/pkg/js"
	"github.com/iwanhae/kubegraph/pkg/watcher"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func main() {
	event, _ := js.NewJSEventEmitter("k8sEvent")

	g, ch := graph.NewGraph()
	go func() {
		for evt := range ch {
			msg := make(map[string]interface{})
			msg["type"] = evt.Status
			msg["id"] = evt.ID
			if node := g.GetNode(evt.ID); node != nil {
				msg["edges"] = node.Edges
				msg["content"] = node.Content
			}
			event.Emit(msg)
			time.Sleep(5 * time.Millisecond)
		}
	}()

	cfg := &rest.Config{
		Host: os.Getenv("HOST"),
	}

	mgr, err := manager.New(cfg, manager.Options{Port: 0})
	if err != nil {
		panic(err)
	}

	for _, o := range []struct {
		obj client.Object
		ctr reconcile.Reconciler
	}{
		{obj: &v1.Namespace{}, ctr: &watcher.NamespaceWatcher{Graph: g}},
		{obj: &v1.Node{}, ctr: &watcher.NodeWatcher{Graph: g}},
		{obj: &v1.Pod{}, ctr: &watcher.PodWatcher{Graph: g}},
		{obj: &v1.Service{}, ctr: &watcher.ServiceWatcher{Graph: g}},
		{obj: &v1.Endpoints{}, ctr: &watcher.EPWatcher{Graph: g}},
		{obj: &appsv1.ReplicaSet{}, ctr: &watcher.ReplicasetWatcher{Graph: g}},
		{obj: &appsv1.Deployment{}, ctr: &watcher.DeploymentWatcher{Graph: g}},
		{obj: &networkingv1.Ingress{}, ctr: &watcher.IngressWatcher{Graph: g}},
	} {
		if err := builder.ControllerManagedBy(mgr).
			Named("kubegraph").
			Watches(&source.Kind{Type: o.obj}, &handler.EnqueueRequestForObject{}).
			Complete(o.ctr); err != nil {
			panic(err)
		}
	}

	if err := mgr.Start(context.Background()); err != nil {
		panic(err)
	}
}
