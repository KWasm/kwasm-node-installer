#!/usr/bin/env sh
set -euo pipefail

cp /assets/crun $NODE_ROOT/usr/local/bin/crun
cp /assets/libwasmedge.so $NODE_ROOT/usr/local/lib/libwasmedge.so
ln -s /usr/local/lib/libwasmedge.so $NODE_ROOT/usr/local/lib/libwasmedge.so.0 && \
ln -s /usr/local/lib/libwasmedge.so $NODE_ROOT/usr/local/lib/libwasmedge.so.0.0.0

if ! grep -q crun $NODE_ROOT/etc/containerd/config.toml; then  
    echo '[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun]
    runtime_type = "io.containerd.runc.v2"
    pod_annotations = ["module.wasm.image/variant", "run.oci.handler"]
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.crun.options]
    BinaryName = "crun"' >> $NODE_ROOT/etc/containerd/config.toml
fi

nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- ldconfig
nsenter -m/$NODE_ROOT/proc/1/ns/mnt -- /bin/systemctl restart containerd
sleep 10