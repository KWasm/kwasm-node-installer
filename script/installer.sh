#!/usr/bin/env sh
set -euo pipefail

KWASM_DIR=/opt/kwasm

CONTAINERD_CONF=/etc/containerd/config.toml
IS_MICROK8S=false
if ps aux | grep kubelet | grep -q snap/microk8s; then
    CONTAINERD_CONF=/var/snap/microk8s/current/args/containerd-template.toml
    IS_MICROK8S=true
    if nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- ls /var/snap/microk8s/current/args/containerd-template.toml > /dev/null 2>&1 ;then
        KWASM_DIR=/var/snap/microk8s/common/kwasm
    else
        echo "Installer seems to run on microk8s but 'containerd-template.toml' not found."
        exit 1
    fi
fi

mkdir -p $NODE_ROOT$KWASM_DIR/bin/
mkdir -p $NODE_ROOT$KWASM_DIR/lib/
case $1 in
    wasmtime)
        cp /assets/crun-wasmtime $NODE_ROOT$KWASM_DIR/bin/crun && \
        cp /assets/libwasmtime.so $NODE_ROOT$KWASM_DIR/lib/libwasmtime.so && \
        ln -sf $KWASM_DIR/lib/libwasmtime.so $NODE_ROOT/lib/libwasmtime.so
        ;;
    *)
    #wasmedge)
        cp /assets/crun-wasmedge $NODE_ROOT$KWASM_DIR/bin/crun && \
        cp /assets/libwasmedge.so $NODE_ROOT$KWASM_DIR/lib/libwasmedge.so && \
        ln -sf $KWASM_DIR/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so && \
        ln -sf $KWASM_DIR/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so.0 && \
        ln -sf $KWASM_DIR/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so.0.0.0
        ;;
esac

cp /assets/containerd-shim-spin-v1 $NODE_ROOT$KWASM_DIR/bin/containerd-shim-spin-v1
cp /assets/containerd-shim-wasmedge-v1 $NODE_ROOT$KWASM_DIR/bin/containerd-shim-wasmedge-v1
if [ -f $NODE_ROOT/usr/local/bin/containerd-shim-spin-v1 ]; then
    # Replace existing spin shim on Azure AKS nodes
    ln -sf $KWASM_DIR/bin/containerd-shim-spin-v1 $NODE_ROOT/usr/local/bin/containerd-shim-spin-v1
    ln -sf $KWASM_DIR/bin/containerd-shim-wasmedge-v1 $NODE_ROOT/usr/local/bin/containerd-shim-spin-v1
elif ! $IS_MICROK8S; then
    ln -sf $KWASM_DIR/bin/containerd-shim-spin-v1 $NODE_ROOT/bin/
    ln -sf $KWASM_DIR/bin/containerd-shim-wasmedge-v1 $NODE_ROOT/bin/
fi

if ! grep -q crun $NODE_ROOT$CONTAINERD_CONF; then  
    echo '[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun]
    runtime_type = "io.containerd.runc.v2"
    pod_annotations = ["module.wasm.image/variant", "run.oci.handler"]
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun.options]
    BinaryName = "'$KWASM_DIR/bin/crun'"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.spin]
    runtime_type = "io.containerd.spin.v1"
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.wasmedge]
    runtime_type = "io.containerd.wasmedge.v1"' >> $NODE_ROOT$CONTAINERD_CONF
    rm -Rf $NODE_ROOT$KWASM_DIR/opt/kwasm/active
fi

if [ ! -f $NODE_ROOT$KWASM_DIR/active ]; then
    if $IS_MICROK8S; then
        nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- systemctl restart snap.microk8s.daemon-containerd
    else
        nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- /bin/systemctl restart containerd
    fi
    touch $NODE_ROOT$KWASM_DIR/active
else
    echo "No change in containerd/config.toml"
fi
