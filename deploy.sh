#!/bin/sh -e

if ! etcdctl get /flitter/domain
then
	echo "No domain set. Please set /flitter/domain."
	exit 1
fi

# deploy flitter
echo "Initiating the container announcer mesh"
fleetctl start -block-attempts 30 ./run/units/flitter-havok.service

echo "Deploying the builder"
fleetctl start -block-attempts 30 ./run/units/flitter-builder.service
fleetctl start -block-attempts 30 ./run/units/flitter-announce@builder.service

echo "Wiring up the builder proxy"
fleetctl start -block-attempts 30 ./run/units/flitter-proxy.service

echo "Spinning up Lagann"
fleetctl start -block-attempts 30 ./run/units/flitter-lagann.service
fleetctl start -block-attempts 30 ./run/units/flitter-announce@lagann.service

echo "Starting the docker registry"
fleetctl start -block-attempts 30 ./run/units/flitter-registry.service
fleetctl start -block-attempts 30 ./run/units/flitter-announce@registry.service

# Routing
echo "Situating the routers"
for n in fleetctl list-machines -no-legend | awk -F'.' '{printf $1"\n"}'
do
	fleetctl start -block-attempts 30 ./run/units/flitter-router@"$n".service
done

echo "Flitter is deployed."
