# fluxmore
Kubernetes Operator that unsuspend flux helmRelease when dependent resources exist in the cluster


### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- helm version 3
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Deploy with Helm**

**Install the controller into the cluster:**

```sh
make docker-build docker-push IMG=<some-registry>/fluxmore:tag
make helm
helm  install mor -n <ns> . --create-namespace
```

**Create a fluxmore instance**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -f config/samples/
```
You can use the `kubectl explain fluxmore`

