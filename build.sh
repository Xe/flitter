#!/bin/sh
# this is a shell script instead of a makefile so it can run in CoreOS

set -e
set -x

# build base image
#docker build -t flitter/init base/image

# build binaries in docker
docker build -t flitter-build:master .

# get binaries out of docker
export CTID=$(docker run -dit flitter-build:master /bin/sh)
docker cp $CTID:/go/bin ./run/
docker rm -f $CTID

# Sorting hat
cd run

# builder
cp bin/builder     builder
cp bin/cloudchaser builder
cp bin/execd       builder

# havok
cp bin/havok       havok

# lagann
cp bin/lagann      lagann

# proxy
cp bin/proxy       proxy

# build active images
for comp in builder havok lagann proxy
do
	docker build -t flitter/$comp:master $comp
done

# clean up
rm builder/builder
rm builder/cloudchaser
rm builder/execd
rm havok/havok
rm lagann/lagann
rm proxy/proxy
