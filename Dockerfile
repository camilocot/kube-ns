FROM golang AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential && \
    apt-get clean && \
    mkdir -p "$GOPATH/src/github.com/camilocot/kube-ns"

ADD . "$GOPATH/src/github.com/camilocot/kube-ns"

RUN cd "$GOPATH/src/github.com/camilocot/kube-ns" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /kubens

FROM alpine:3.12

COPY --from=builder /kubens /bin/kubens

ENTRYPOINT ["/bin/kubens"]
