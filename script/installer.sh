#!/usr/bin/env sh
set -euo pipefail

KWASM_DIR=/opt/kwasm

CONTAINERD_CONF=/etc/containerd/config.toml
IS_MICROK8S=false
IS_K3S=false
if ps aux | grep kubelet | grep -q snap/microk8s; then
    CONTAINERD_CONF=/var/snap/microk8s/current/args/containerd-template.toml
    IS_MICROK8S=true
    if nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- ls /var/snap/microk8s/current/args/containerd-template.toml > /dev/null 2>&1 ;then
        KWASM_DIR=/var/snap/microk8s/common/kwasm
    else
        echo "Installer seems to run on microk8s but 'containerd-template.toml' not found."
        exit 1
    fi
elif ls $NODE_ROOT/var/lib/rancher/k3s/agent/etc/containerd/config.toml > /dev/null 2>&1 ; then
    IS_K3S=true
    cp $NODE_ROOT/var/lib/rancher/k3s/agent/etc/containerd/config.toml $NODE_ROOT/var/lib/rancher/k3s/agent/etc/containerd/config.toml.tmpl
    CONTAINERD_CONF=/var/lib/rancher/k3s/agent/etc/containerd/config.toml.tmpl
fi

IS_ALPINE=false
CRUN_WASMEDGE=crun-wasmedge
LIB_WASMEDGE=libwasmedge.so
if grep -iq alpine $NODE_ROOT/etc/issue 2>/dev/null ; then
    IS_ALPINE=true
    CRUN_WASMEDGE=crun-wasmedge-musl
    LIB_WASMEDGE=libwasmedge-musl.so
    nsenter --target 1 --mount --uts --ipc --net -- sh -c "which apk && apk add libseccomp lld-libs"
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
        cp /assets/$CRUN_WASMEDGE $NODE_ROOT$KWASM_DIR/bin/crun && \
        cp /assets/$LIB_WASMEDGE $NODE_ROOT$KWASM_DIR/lib/libwasmedge.so && \
        ln -sf $KWASM_DIR/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so && \
        ln -sf $KWASM_DIR/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so.0 && \
        ln -sf $KWASM_DIR/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so.0.0.0
        ;;
esac

cp /assets/containerd-shim-spin-v1 $NODE_ROOT$KWASM_DIR/bin/containerd-shim-spin-v1
cp /assets/containerd-shim-wasmedge-v1 $NODE_ROOT$KWASM_DIR/bin/containerd-shim-wasmedge-v1
cp /assets/containerd-shim-wws-v1 $NODE_ROOT$KWASM_DIR/bin/containerd-shim-wws-v1
if [ -f $NODE_ROOT/usr/local/bin/containerd-shim-spin-v1 ]; then
    # Replace existing spin shim on Azure AKS nodes
    ln -sf $KWASM_DIR/bin/containerd-shim-spin-v1 $NODE_ROOT/usr/local/bin/containerd-shim-spin-v1
    ln -sf $KWASM_DIR/bin/containerd-shim-wasmedge-v1 $NODE_ROOT/usr/local/bin/containerd-shim-wasmedge-v1
    ln -sf $KWASM_DIR/bin/containerd-shim-wws-v1 $NODE_ROOT/usr/local/bin/containerd-shim-wws-v1
elif ! $IS_MICROK8S; then
    ln -sf $KWASM_DIR/bin/containerd-shim-spin-v1 $NODE_ROOT/bin/
    ln -sf $KWASM_DIR/bin/containerd-shim-wasmedge-v1 $NODE_ROOT/bin/
    ln -sf $KWASM_DIR/bin/containerd-shim-wws-v1 $NODE_ROOT/bin/
fi

CRI='"io.containerd.grpc.v1.cri"'
if $IS_K3S; then
    CRI='cri'
fi
if ! grep -q crun $NODE_ROOT$CONTAINERD_CONF; then
    echo '[plugins.'$CRI'.containerd.runtimes.crun]
    runtime_type = "io.containerd.runc.v2"
    pod_annotations = ["module.wasm.image/variant", "run.oci.handler"]
[plugins.'$CRI'.containerd.runtimes.crun.options]
    BinaryName = "'$KWASM_DIR/bin/crun'"
[plugins.'$CRI'.containerd.runtimes.spin]
    runtime_type = "io.containerd.spin.v1"
[plugins.'$CRI'.containerd.runtimes.wasmedge]
    runtime_type = "io.containerd.wasmedge.v1"
[plugins.'$CRI'.containerd.runtimes.wws]
    runtime_type = "io.containerd.wws.v1"' >> $NODE_ROOT$CONTAINERD_CONF
    rm -Rf $NODE_ROOT$KWASM_DIR/opt/kwasm/active
fi

if [ ! -f $NODE_ROOT$KWASM_DIR/active ]; then
    touch $NODE_ROOT$KWASM_DIR/active
    if $IS_MICROK8S; then
        nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- systemctl restart snap.microk8s.daemon-containerd
    elif ls $NODE_ROOT/etc/init.d/containerd > /dev/null 2>&1 ; then
        nsenter --target 1 --mount --uts --ipc --net -- /etc/init.d/containerd restart
    elif ls $NODE_ROOT/etc/init.d/k3s > /dev/null 2>&1 ; then
        nsenter --target 1 --mount --uts --ipc --net -- /etc/init.d/k3s restart
    else
        nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- /bin/systemctl restart containerd
    fi
else
    echo "No change in containerd/config.toml"
fi
