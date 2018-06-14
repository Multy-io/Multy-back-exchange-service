# docker build --force-rm -t multy-back-exchange-info:latest .
FROM golang:1.9.4

ENV SRC_BASE=github.com/Enmk/Multy-back-exchange-service
RUN go get ${SRC_BASE} && \
    cd $GOPATH/src/${SRC_BASE} && \
    make all-with-deps dist

WORKDIR /go/src/${SRC_BASE}/cmd

ENTRYPOINT $GOPATH/src/${SRC_BASE}/cmd/exchanger
