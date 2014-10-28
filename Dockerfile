FROM xena/golang
MAINTAINER Xena <xena@yolo-swag.com>

ENV CGO_ENABLED 0

RUN go get -v -a -ldflags '-s' github.com/Xe/flitter/...
