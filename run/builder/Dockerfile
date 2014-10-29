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

RUN chmod a+x /usr/bin/docker

ADD . /app
ADD runit/ /etc/service
