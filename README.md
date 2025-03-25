# deck8sops-dev

`deck8sops-dev` is a tool designed to declaratively manage Kubernetes manifest operations on a [kind](https://github.com/kubernetes-sigs/kind) based Kubernetes cluster. Define your required operations and their configurations in a YAML file and apply them with a single command.

## Features

- Declarative manifest operation definitions using YAML
- Support for both Helm charts and Kubectl manifests
- Sequential installation and reverse-order uninstallation

## Installation

```bash
$ go build -o deck8sops-dev cmd/main.go
```

## Usage

### Launch kubernetes

```bash
cat << EOF > test-cluster.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: test
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
EOF
```

```bash
$ kind create cluster --config test-cluster.yaml
```

### Creating a configuration file

See `examples` directory for sample configurations. You can define both Helm and kubectl operations.

### Installing operations

```bash
$ deck8sops-dev create -config examples/kyverno/operations.yaml
```

### Uninstalling operations

```bash
$ deck8sops-dev delete -config examples/kyverno/operations.yaml
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
