#!/bin/sh -ex

if ! etcdctl get /flitter/domain
then
	echo "No domain set. Please set /flitter/domain."
	exit 1
fi

function start-fleet-and-wait {
  fleetctl load $*

	for unit in $*
	do
		fleetctl ssh $unit sudo systemctl start $unit
	done
}

# deploy flitter
echo "Initiating the container announcer mesh"
fleetctl start -block-attempts 30 ./run/units/flitter-havok.service

echo "Deploying the builder"
start-fleet-and-wait ./run/units/flitter-builder.service
start-fleet-and-wait ./run/units/flitter-announce@builder.service

echo "Wiring up the builder proxy"
fleetctl start -block-attempts 30 ./run/units/flitter-proxy.service

echo "Spinning up Lagann"
start-fleet-and-wait ./run/units/flitter-lagann.service
start-fleet-and-wait ./run/units/flitter-announce@lagann.service

echo "Starting the docker registry"
start-fleet-and-wait ./run/units/flitter-registry.service
start-fleet-and-wait ./run/units/flitter-announce@registry.service

# Routing
echo "Situating the routers"
for n in $(fleetctl list-machines -no-legend | awk -F'.' '{printf $1"\n"}')
do
	start-fleet-and-wait ./run/units/flitter-router@"$n".service
done

echo "Flitter is deployed."
