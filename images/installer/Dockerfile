FROM ubuntu:18.04 as builder-crun

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y curl make git gcc build-essential pkgconf libtool libsystemd-dev libprotobuf-c-dev libcap-dev libseccomp-dev libyajl-dev go-md2man libtool autoconf python3 automake xz-utils \
    && curl https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash -s -- -p /usr/local --version=0.11.0 \
    && git clone --depth 1 --branch 1.6 https://github.com/containers/crun.git \
    && cd crun \
    && ./autogen.sh \
    && ./configure --with-wasmedge --enable-embedded-yajl\
    && make \
    && mv crun crun-wasmedge
RUN curl -L https://github.com/bytecodealliance/wasmtime/releases/download/v1.0.1/wasmtime-v1.0.1-$(uname -m)-linux-c-api.tar.xz | tar xJf - -C / \
    && cp -R /wasmtime-v1.0.1-$(uname -m)-linux-c-api/* /usr/local/ \
    && cd /crun \
    && ./configure --with-wasmtime --enable-embedded-yajl \
    && make \
    && mv crun crun-wasmtime

FROM busybox as crun-wasmtime

COPY --from=builder-crun /crun/crun-wasmtime /assets/crun-wasmtime
COPY --from=builder-crun /usr/local/lib/libwasmtime.so /assets/libwasmtime.so
COPY script/installer.sh /script/installer.sh

CMD sh /script/installer.sh wasmtime

FROM busybox as crun-wasmedge

COPY --from=builder-crun /crun/crun-wasmedge /assets/crun-wasmedge
COPY --from=builder-crun /usr/local/lib/libwasmedge.so /assets/libwasmedge.so
COPY script/installer.sh /script/installer.sh

CMD sh /script/installer.sh wasmedge

FROM busybox

COPY --from=builder-crun /crun/crun-* /assets/
COPY --from=builder-crun /usr/local/lib/libwasmedge.so /assets/libwasmedge.so
COPY --from=builder-crun /usr/local/lib/libwasmtime.so /assets/libwasmtime.so
COPY script/installer.sh /script/installer.sh

CMD sh /script/installer.sh wasmedge