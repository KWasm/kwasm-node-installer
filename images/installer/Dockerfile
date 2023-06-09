FROM ubuntu:18.04 as builder-crun

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y curl make git gcc build-essential pkgconf libtool libsystemd-dev libprotobuf-c-dev libcap-dev libseccomp-dev libyajl-dev go-md2man libtool autoconf python3 automake xz-utils \
    && curl https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash -s -- -p /usr/local --version=0.11.2 \
    && git clone --depth 1 --branch 1.8.1 https://github.com/containers/crun.git \
    && cd crun \
    && ./autogen.sh \
    && ./configure --with-wasmedge --enable-embedded-yajl \
    && make \
    && mv crun crun-wasmedge
RUN curl -L https://github.com/bytecodealliance/wasmtime/releases/download/v1.0.1/wasmtime-v1.0.1-$(uname -m)-linux-c-api.tar.xz | tar xJf - -C / \
    && cp -R /wasmtime-v1.0.1-$(uname -m)-linux-c-api/* /usr/local/ \
    && cd /crun \
    && ./configure --with-wasmtime --enable-embedded-yajl \
    && make \
    && mv crun crun-wasmtime

FROM alpine:3.18 as builder-crun-alpine
RUN apk add gcc automake autoconf libtool gettext pkgconf make musl-dev python3 libcap-dev libseccomp-dev yajl-dev argp-standalone go-md2man curl bash coreutils cmake ninja clang lld-dev llvm-dev boost-dev llvm-static git \
    && cd / && git clone --depth 1 --branch 0.12.1 https://github.com/WasmEdge/WasmEdge.git \
    && cd WasmEdge && bash -c "export CC=clang && export CXX=clang++ && cmake -Bbuildmusl -GNinja -DCMAKE_BUILD_TYPE=Release . && cmake --build buildmusl" \
    && curl https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash -s -- -p /usr/local --version=0.12.1 \
    && cd / && git clone --depth 1 --branch 1.8.1 https://github.com/containers/crun.git \
    && cd crun \
    && ./autogen.sh \    
    && ./configure --with-wasmedge --enable-embedded-yajl --disable-systemd \
    && make \
    && mv crun crun-wasmedge-musl

FROM ubuntu:18.04 as builder-containerd-shim

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y curl

RUN curl -L https://github.com/deislabs/containerd-wasm-shims/releases/download/v0.6.0/containerd-wasm-shims-v1-linux-$(uname -m).tar.gz | tar -xzf - \
    && curl -L https://github.com/second-state/runwasi/releases/download/v0.3.2/containerd-shim-wasmedge-v1-v0.3.2-linux-$(uname -m | sed s/aarch64/arm64/g | sed s/x86_64/amd64/g).tar.gz | tar -xzf -

FROM ubuntu:18.04 as builder-containerd-shim-spin
RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y git curl libclang-dev make build-essential clang\
    && git clone --depth 1 --branch main https://github.com/deislabs/containerd-wasm-shims.git
WORKDIR /containerd-wasm-shims
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > rustup && chmod +x ./rustup \
    && ./rustup -y \
    && bash -c "source ~/.cargo/env && make build-spin"

FROM busybox as crun-wasmtime

COPY --link --from=builder-crun /crun/crun-wasmtime /assets/crun-wasmtime
COPY --link --from=builder-crun /usr/local/lib/libwasmtime.so /assets/libwasmtime.so
COPY script/installer.sh /script/installer.sh

CMD sh /script/installer.sh wasmtime

FROM busybox as crun-wasmedge

COPY --link --from=builder-crun /crun/crun-wasmedge /assets/crun-wasmedge
COPY --link --from=builder-crun-alpine /crun/crun-wasmedge-musl /assets/crun-wasmedge-musl
COPY --link --from=builder-crun /usr/local/lib/libwasmedge.so /assets/libwasmedge.so
COPY script/installer.sh /script/installer.sh

CMD sh /script/installer.sh wasmedge

FROM busybox

COPY script/installer.sh /script/installer.sh
COPY --link --from=builder-crun /crun/crun-* /assets/
COPY --link --from=builder-crun /usr/local/lib/libwasmedge.so /assets/libwasmedge.so
COPY --link --from=builder-crun /usr/local/lib/libwasmtime.so /assets/libwasmtime.so
COPY --link --from=builder-crun-alpine /crun/crun-wasmedge-musl /assets/crun-wasmedge-musl
COPY --link --from=builder-crun-alpine /WasmEdge/buildmusl/lib/api/libwasmedge.so /assets/libwasmedge-musl.so
COPY --link --from=builder-containerd-shim /containerd-shim-slight-v1 /assets/containerd-shim-slight-v1
COPY --link --from=builder-containerd-shim /containerd-shim-spin-v1 /assets/containerd-shim-spin-v1
COPY --link --from=builder-containerd-shim /containerd-shim-wasmedge-v1 /assets/containerd-shim-wasmedge-v1
#COPY --link --from=builder-containerd-shim-spin /containerd-wasm-shims/containerd-shim-spin-v1/target/release/containerd-shim-spin-v1 /assets/containerd-shim-spin-v1

CMD sh /script/installer.sh wasmedge
