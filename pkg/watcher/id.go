package watcher

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NodeID(name string) string {
	return fmt.Sprintf("node/%s", name)
}
func NSID(namespace string) string {
	return fmt.Sprintf("ns/%s", namespace)
}
func DeploymentID(namespace string, name string) string {
	return fmt.Sprintf("deploy/%s/%s", namespace, name)
}
func IngressID(namespace string, name string) string {
	return fmt.Sprintf("ing/%s/%s", namespace, name)
}

func ReplicaSetID(namespace string, name string) string {
	return fmt.Sprintf("rs/%s/%s", namespace, name)
}
func PodID(namespace string, name string) string {
	return fmt.Sprintf("pod/%s/%s", namespace, name)
}

func EPID(namespace string, name string) string {
	return fmt.Sprintf("endpoint/%s/%s", namespace, name)
}

func SvcID(namespace string, name string) string {
	return fmt.Sprintf("service/%s/%s", namespace, name)
}

func IP(ip string) string {
	return fmt.Sprintf("ip/%s", ip)
}

func GVK(gvk metav1.TypeMeta) string {
	return fmt.Sprintf("%s/%s", gvk.APIVersion, gvk.Kind)
}

type Content struct {
	Color     string `json:"color"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func HasEdge(edges []string, id string) (idx int) {
	for idx := 0; idx <= len(edges); idx++ {
		if edges[idx] == id {
			return idx
		}
	}
	return -1
}
