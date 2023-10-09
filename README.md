# KWasm-node-installer

> A Go version of this installer is developed at the [go-rewrite](https://github.com/KWasm/kwasm-node-installer/tree/go-rewrite) branch.

The KWasm Node Installer is a container image that contains binaries and configuration to enable a Kubernetes node to run pure webassembly images.

> WARNING: Not meant to be used in production!

Since this installer changes the configuration of the node it can make a node unusable. We recommend using a fresh KinD/MiniKube/MicroK8s or a managed cloud service like AKS/GKE/EKS.

## Supported Kubernetes distributions

- KinD
- MiniKube
- MicroK8s
- Rancher RKE2
- Azure AKS
- GCP GKE (Ubuntu Nodes)
- AWS EKS (AmazonLinux2)
- AWS EKS (Ubuntu Nodes)
- Digital Ocean Kubernetes

## Currently not supported Kubernetes distributions

- OCI OKE
- OpenShift

## Quickstart

### KinD

Prerequisites:

- Docker
- KinD
- kubectl

```bash
kind create cluster

# As default crun-wasmedge is used for installation.
kubectl apply -f https://raw.githubusercontent.com/KWasm/kwasm-node-installer/main/example/daemonset.yaml


kubectl apply -f https://raw.githubusercontent.com/KWasm/kwasm-node-installer/main/example/test-job.yaml
```

## How to build

The dockerfile `images/installer/Dockerfile` has multiple stages by building it you create an image that contains a crun version built with WasmEdge support and one with WasmTime. If you only need one of them you can use a different stage with the target parameter `--target crun-wasmedge` or `--target crun-wasmtime`.

```bash
docker build . -f images/installer/Dockerfile -t kwasm/kwasm-installer
```
