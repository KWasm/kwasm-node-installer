ARG CONTAINERD_RUNWASI="download-containerd-runwasi"

FROM ubuntu:22.04 AS download-containerd-runwasi
RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y curl

RUN mkdir -p /release/bin/ \
    && curl -L https://github.com/containerd/runwasi/releases/download/containerd-shim-wasmedge%2Fv0.3.0/containerd-shim-wasmedge-$(uname -m | sed s/arm64/aarch64/g | sed s/amd64/x86_64/g).tar.gz | tar -xzf - -C /release/bin/ \
    && curl -L https://github.com/containerd/runwasi/releases/download/containerd-shim-wasmtime%2Fv0.3.0/containerd-shim-wasmtime-$(uname -m | sed s/arm64/aarch64/g | sed s/amd64/x86_64/g).tar.gz | tar -xzf - -C /release/bin/ \
    && curl -L https://github.com/containerd/runwasi/releases/download/containerd-shim-wasmer%2Fv0.3.0/containerd-shim-wasmer-$(uname -m | sed s/arm64/aarch64/g | sed s/amd64/x86_64/g).tar.gz | tar -xzf - -C /release/bin/

FROM ${CONTAINERD_RUNWASI} AS containerd_runwasi

FROM ubuntu:22.04 AS deislabs_containerd-wasm-shims

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y curl

RUN mkdir /assets \
    && curl -L https://github.com/deislabs/containerd-wasm-shims/releases/download/v0.9.2/containerd-wasm-shims-v1-lunatic-linux-$(uname -m).tar.gz  | tar -xzf - -C /assets \
    && curl -L https://github.com/deislabs/containerd-wasm-shims/releases/download/v0.9.2/containerd-wasm-shims-v1-slight-linux-$(uname -m).tar.gz  | tar -xzf - -C /assets \
    && curl -L https://github.com/deislabs/containerd-wasm-shims/releases/download/v0.9.2/containerd-wasm-shims-v1-spin-linux-$(uname -m).tar.gz  | tar -xzf - -C /assets \
    && curl -L https://github.com/deislabs/containerd-wasm-shims/releases/download/v0.9.2/containerd-wasm-shims-v1-wws-linux-$(uname -m).tar.gz  | tar -xzf - -C /assets


FROM busybox

COPY script/installer.sh /script/installer.sh
COPY --link --from=deislabs_containerd-wasm-shims /assets /assets
COPY --link --from=containerd_runwasi /release/bin/containerd-shim-wasmedge-v1 /assets/
COPY --link --from=containerd_runwasi /release/bin/containerd-shim-wasmer-v1 /assets/
COPY --link --from=containerd_runwasi /release/bin/containerd-shim-wasmtime-v1 /assets/
CMD sh /script/installer.sh wasmedge
