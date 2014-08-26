FROM deis/go

ADD . /go/src/github.com/Xe/flitter

RUN go get github.com/Xe/flitter/builder && \
    go get github.com/Xe/flitter/execd && \
    go get github.com/Xe/flitter/cloudchaser
