#!/bin/sh -ex

if ! etcdctl --peers http://$FLEETCTL_TUNNEL:4001 get /flitter/domain
then
	echo "No domain set. Please set /flitter/domain."
	exit 1
fi

function fleet_start_wait {
  fleetctl load $*

	for unit in $*
	do
		fleetctl ssh $unit sudo systemctl start $unit
	done
}

cd ./run/units/

# deploy flitter
echo "Initiating the container announcer mesh"
fleetctl start -block-attempts 30 flitter-havok.service

echo "Deploying the builder"
fleet_start_wait flitter-builder.service
fleet_start_wait flitter-announce@builder.service

echo "Wiring up the builder proxy"
fleetctl start -block-attempts 30 flitter-proxy.service

echo "Spinning up Lagann"
fleet_start_wait flitter-lagann.service
fleet_start_wait flitter-announce@lagann.service

echo "Starting the docker registry"
fleet_start_wait flitter-registry.service
fleet_start_wait flitter-announce@registry.service

# Routing
echo "Situating the routers"
for n in $(fleetctl list-machines -no-legend | awk -F'.' '{printf $1"\n"}')
do
	fleet_start_wait flitter-router@"$n".service
done

echo "Flitter is deployed."
