FROM flitter/init
MAINTAINER Xena <xena@yolo-swag.com>

ENV GITHOME /home/git
ENV GITUSER git

# Configure base builder image
RUN apt-get update && apt-get install git -y
RUN useradd -d $GITHOME $GITUSER && mkdir $GITHOME && chown -R $GITUSER:$GITUSER $GITHOME && \
    echo "%git    ALL=(ALL:ALL) NOPASSWD:/app/bin/builder" >> /etc/sudoers

ADD https://get.docker.com/builds/Linux/x86_64/docker-latest /usr/bin/docker
RUN chmod a+x /usr/bin/docker

VOLUME /var/lib/docker
WORKDIR /app
EXPOSE 22

ADD . /app
ADD runit/ /etc/service
