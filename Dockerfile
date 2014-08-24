FROM deis/go

ADD . /go/src/github.com/Xe/flitter
ENV CGO_ENABLED 0

RUN go get -a -ldflags '-s' github.com/Xe/flitter/builder && \
    go get -a -ldflags '-s' github.com/Xe/flitter/execd
