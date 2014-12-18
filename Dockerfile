FROM xena/golang
MAINTAINER Xena <xena@yolo-swag.com>

ENV CGO_ENABLED 0

ADD . /go/src/github.com/Xe/flitter

RUN go get github.com/tools/godep
RUN godep go get -v -a -ldflags '-s' github.com/Xe/flitter/...
