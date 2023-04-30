# kubegraph

Visualizing Kubernetes resources as a graph. It uses [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) to get Kubernetes resources in realtime from your web browser and [d3-force](https://github.com/d3/d3-force) to render graph.

[![Video Label](http://img.youtube.com/vi/5zH78ZByePc/0.jpg)](https://youtu.be/5zH78ZByePc)

Not neccessary, but for ease of use (there's a bunch of problems related to authentication and CORS to make request to your `kube-apiserver` in web browser), this project relies on `kubectl proxy` command.

## How to use

```bash
curl -L https://github.com/iwanhae/kubegraph/releases/download/v0.0.1/kubegraph.tar -o kubegraph.tar
tar -xvf kubegraph.tar
kubectl proxy -w ./static

# open
# http://localhost:8001/static
# in your web browser
```
