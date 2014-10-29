#!/bin/sh

export TAG={$1:-latest}

for comp in builder havok lagann proxy;
do
	docker tag flitter/$comp:master flitter/$comp:latest
	docker push flitter/$comp:latest
done
