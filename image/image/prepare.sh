#!/bin/bash
set -e
source /build/buildconfig
set -x

## Temporarily disable dpkg fsync to make building faster.
echo force-unsafe-io > /etc/dpkg/dpkg.cfg.d/02apt-speedup

## Prevent initramfs updates from trying to run grub and lilo.
## https://journal.paul.querna.org/articles/2013/10/15/docker-ubuntu-on-rackspace/
## http://bugs.debian.org/cgi-bin/bugreport.cgi?bug=594189
export INITRD=no
mkdir -p /etc/container_environment
echo -n no > /etc/container_environment/INITRD

## Enable Ubuntu Universe and Multiverse.
sed -i 's/^#\s*\(deb.*universe\)$/\1/g' /etc/apt/sources.list
sed -i 's/^#\s*\(deb.*multiverse\)$/\1/g' /etc/apt/sources.list
apt-get update

## Fix some issues with APT packages.
## See https://github.com/dotcloud/docker/issues/1024
dpkg-divert --local --rename --add /sbin/initctl
ln -sf /bin/true /sbin/initctl

## Replace the 'ischroot' tool to make it always return true.
## Prevent initscripts updates from breaking /dev/shm.
## https://journal.paul.querna.org/articles/2013/10/15/docker-ubuntu-on-rackspace/
## https://bugs.launchpad.net/launchpad/+bug/974584
dpkg-divert --local --rename --add /usr/bin/ischroot
ln -sf /bin/true /usr/bin/ischroot

## Workaround https://github.com/dotcloud/docker/issues/2267,
## not being able to modify /etc/hosts.
mkdir -p /etc/workaround-docker-2267
ln -s /etc/workaround-docker-2267 /cte
cp /build/bin/workaround-docker-2267 /usr/bin/

## Install HTTPS support for APT.
$minimal_apt_get_install apt-transport-https ca-certificates

## Install add-apt-repository
$minimal_apt_get_install software-properties-common

## Upgrade all packages.
apt-get dist-upgrade -y --no-install-recommends

## Fix locale.
$minimal_apt_get_install language-pack-en
locale-gen en_US.UTF-8

## teh deps
$minimal_apt_get_install make net-tools sudo wget vim strace lsof netcat

## Download our version of etcdctl
wget -q https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.5 -O /usr/local/bin/etcdctl
chmod +x /usr/local/bin/etcdctl

## Download confd
wget -q https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json -O /usr/local/bin/confd
chmod +x /usr/local/bin/confd
