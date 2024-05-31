FROM golang:1.17.3-alpine
RUN apk update && apk add \
    curl \
    bind-tools \
    iptables \
    jq \
    nghttp2 \
    tcpdump \
    iproute2 \
    lsof \
    conntrack-tools \
    tshark && \
    rm -rf /var/lib/apt/lists/*

# We still rely on old iptables-legacy syntax.
#RUN update-alternatives --set iptables /usr/sbin/iptables-legacy \
#    && update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy

COPY websocket /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build \
     -o /out/echo-client -tags prod -mod=readonly -ldflags "-s -w" \
     /build/examples/echo/client.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build \
    -o /out/echo-server -tags prod -mod=readonly -ldflags "-s -w" \
    /build/examples/echo/server.go

ENTRYPOINT [ "tshark", "-i", "any" ]