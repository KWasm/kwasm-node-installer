# kwasm-node-installer
Installs KWasm on Kubernetes nodes.

## How to build
```bash
docker build . -f images/installer/Dockerfile -t kwasm/kwasm-installer
```

## Quickstart
```bash
kind create cluster

kubectl apply -f https://raw.githubusercontent.com/KWasm/kwasm-node-installer/main/example/daemonset.yaml

echo 'apiVersion: node.k8s.io/v1                             
kind: RuntimeClass
metadata:
  name: crun
handler: crun' |kubectl apply -f -

kubectl apply -f https://raw.githubusercontent.com/KWasm/kwasm-node-installer/main/example/test-job.yaml
```
