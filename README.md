# graph

Visualize your [Knative Eventing](http://github.com/knative/eventing)
connections.

## Usage

Visit the root of the graph service in a web browser. This will show you the
graph of the current Knative resources in the namespace the graph resource is
installed.

> Note: Work is required to support installation of `graph` into multiple
> namespaces.

# Deploying

## From Release v0.2.0

To install into the `default` namespace,

```shell
kubectl apply -f https://github.com/n3wscott/graph/releases/download/v0.2.0/graph.yaml
```

To install into a `test` namespace,

```shell
export NAMESPACE=test # <-- update test to your target namespace.
curl -L https://github.com/n3wscott/graph/releases/download/v0.2.0/graph.yaml \
  | sed "s/default/${NAMESPACE}/" \
  | kubectl apply -n $NAMESPACE --filename -
```

## From Source

To install into the `default` namespace,

```shell
ko apply -f config/
```

To install into a `test` namespace,

```shell
export NAMESPACE=test # <-- update test to your target namespace.
ko resolve -f config/ \
  | sed "s/default/${NAMESPACE}/" \
  | kubectl apply -n $NAMESPACE --filename -
```
