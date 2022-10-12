#!/usr/bin/env sh
set -euo pipefail

mkdir -p $NODE_ROOT/opt/kwasm/bin/
mkdir -p $NODE_ROOT/opt/kwasm/lib/
case $1 in
    wasmtime)
        cp /assets/crun-wasmtime $NODE_ROOT/opt/kwasm/bin/crun && \
        cp /assets/libwasmtime.so $NODE_ROOT/opt/kwasm/lib/libwasmtime.so && \
        ln -sf /opt/kwasm/lib/libwasmtime.so $NODE_ROOT/lib/libwasmtime.so
        ;;
    *)
    #wasmedge)
        cp /assets/crun-wasmedge $NODE_ROOT/opt/kwasm/bin/crun && \
        cp /assets/libwasmedge.so $NODE_ROOT/opt/kwasm/lib/libwasmedge.so && \
        ln -sf /opt/kwasm/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so && \
        ln -sf /opt/kwasm/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so.0 && \
        ln -sf /opt/kwasm/lib/libwasmedge.so $NODE_ROOT/lib/libwasmedge.so.0.0.0
        ;;

esac

if ! grep -q crun $NODE_ROOT/etc/containerd/config.toml; then  
    echo '[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun]
    runtime_type = "io.containerd.runc.v2"
    pod_annotations = ["module.wasm.image/variant", "run.oci.handler"]
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun.options]
    BinaryName = "/opt/kwasm/bin/crun"' >> $NODE_ROOT/etc/containerd/config.toml
    rm -Rf $NODE_ROOT/opt/kwasm/active
fi

if [ ! -f $NODE_ROOT/opt/kwasm/active ]; then
    nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- /bin/systemctl restart containerd
    touch $NODE_ROOT/opt/kwasm/active
else
    echo "No change in containerd/config.toml"
fi
