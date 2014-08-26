FROM deis/go

ADD . /go/src/github.com/Xe/flitter

# git2go is a special snowflake with its installation
# The first `go get` call will fail.
RUN apt-get update && apt-get install -y cmake && \
    go get github.com/libgit2/git2go ; cd $GOPATH/src/github.com/libgit2/git2go && \
    git submodule update --init && make install && \
    go get github.com/Xe/flitter/builder && \
    go get github.com/Xe/flitter/execd
