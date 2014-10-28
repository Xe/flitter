FROM xena/golang
MAINTAINER Xena <xena@yolo-swag.com>

ADD . /go/src/github.com/Xe/flitter

ENV CGO_ENABLED 0

RUN go get -v -a -ldflags '-s' github.com/Xe/flitter/builder/... && \
    go get -v -a -ldflags '-s' github.com/Xe/flitter/execd/... && \
    go get -v -a -ldflags '-s' github.com/Xe/flitter/cloudchaser/...
