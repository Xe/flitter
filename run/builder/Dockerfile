FROM flitter/init
MAINTAINER Xena <xena@yolo-swag.com>

ENV GITHOME /home/git
ENV GITUSER git

# Configure base builder image
RUN apt-get update && apt-get install git -y
RUN wget -qO- https://get.docker.com/ | sh
RUN useradd -d $GITHOME $GITUSER && mkdir $GITHOME && chown -R $GITUSER:$GITUSER $GITHOME && \
    echo "%git    ALL=(ALL:ALL) NOPASSWD:/app/bin/builder" >> /etc/sudoers

VOLUME /var/lib/docker
WORKDIR /app
EXPOSE 22

ADD https://get.docker.com/builds/Linux/x86_64/docker-1.2.0 /usr/bin/docker
RUN chmod a+x /usr/bin/docker

ADD https://github.com/coreos/fleet/releases/download/v0.8.3/fleet-v0.8.3-linux-amd64.tar.gz /opt/fleet
ENV PATH $PATH:/opt/fleet

ADD . /app
ADD runit/ /etc/service

CMD /app/dind
