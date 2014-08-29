FROM flitter/init

# Set up environment variables.
ENV PATH /usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go
ENV GOPATH /home/goworld
ENV VULCANPATH /home/goworld/src/github.com/mailgun/vulcand

# Let's install go just like Docker (from source).
RUN apt-get update -q && \
    apt-get install -qy build-essential curl git && \
    curl -s https://storage.googleapis.com/golang/go1.3.1.src.tar.gz | tar -v -C /usr/local -xz && \
    cd /usr/local/go/src && ./make.bash --no-clean 2>&1 && \
    cd $VULCANPATH && make install && \
    mkdir /opt/vulcan && \
    cp /home/goworld/bin/vulcand /opt/vulcan && \
    cp /home/goworld/bin/vulcanctl /opt/vulcan && \
    echo "Cleanup" && \
    rm -rf /usr/local/go /home/goworld

ADD run /etc/service/vulcand/run

ENTRYPOINT /sbin/my_init
