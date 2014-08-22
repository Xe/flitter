FROM deis/init

ENV GOPATH /go

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y git golang mercurial && mkdir /go && \
    go get github.com/Xe/flitter/builder && go get github.com/Xe/flitter/execd

ENTRYPOINT /sbin/my_init
