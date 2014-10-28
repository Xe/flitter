FROM flitter/init

RUN apt-get update && apt-get install -y golang git-core mercurial --no-install-recommends

ENV GOPATH /go

ADD . /go/src/github.com/Xe/lagann

RUN go get github.com/Xe/lagann/...

ADD run/init /etc/service/lagann/run

EXPOSE 3000

CMD /sbin/my_init
