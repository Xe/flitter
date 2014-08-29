FROM flitter/init
MAINTAINER Xena <xena@yolo-swag.com>

ENV GITHOME /home/git
ENV GITUSER git

# Install docker in docker and configure base builder image
RUN echo "deb http://get.docker.io/ubuntu docker main" > /etc/apt/sources.list.d/docker.list && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9 && \
    apt-get update && apt-get install -yq openssh-server git aufs-tools iptables lxc curl lxc-docker-1.1.2 && \
    ln -s /usr/share/i18n/SUPPORTED /var/lib/locales/supported.d/all && locale-gen && \
    useradd -d $GITHOME $GITUSER && mkdir $GITHOME && chown -R $GITUSER:$GITUSER $GITHOME && \
    echo "%git    ALL=(ALL:ALL) NOPASSWD:/app/bin/builder" >> /etc/sudoers

VOLUME /var/lib/docker
WORKDIR /app
EXPOSE 22

ADD . /app
ADD runit/ /etc/service
