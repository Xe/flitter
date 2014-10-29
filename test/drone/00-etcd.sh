#!/bin/bash

set -ex

# download etcd
mkdir ~/tmp || true
cd ~/tmp
wget https://github.com/coreos/etcd/releases/download/v0.4.6/etcd-v0.4.6-linux-amd64.tar.gz > /dev/null
tar xf etcd-v0.4.6-linux-amd64.tar.gz
etcd-v0.4.6-linux-amd64/etcd -data-dir $(uuidgen) &

sleep 5
