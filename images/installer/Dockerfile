FROM ubuntu:22.04 as builder-wasmedge

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt install -y curl make git gcc build-essential pkgconf libtool libsystemd-dev libprotobuf-c-dev libcap-dev libseccomp-dev libyajl-dev go-md2man libtool autoconf python3 automake \
    && curl https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash -s -- -p /usr/local --version=0.11.0 \
    && git clone --depth 1 --branch 1.6 https://github.com/containers/crun.git \
    && cd crun \
    && ./autogen.sh \
    && ./configure --with-wasmedge --enable-embedded-yajl\
    && make

FROM busybox

COPY --from=builder-wasmedge /crun/crun /assets/crun
COPY --from=builder-wasmedge /usr/local/lib/libwasmedge.so /assets/libwasmedge.so
COPY script/installer.sh /script/installer.sh

CMD sh /script/installer.sh