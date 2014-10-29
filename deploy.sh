#!/bin/sh

set -ex

# deploy flitter
fleetctl start ./run/units/flitter-registry.service
fleetctl start ./run/units/flitter-registry-announce.service

fleetctl start ./run/units/flitter-lagann.service
fleetctl start ./run/units/flitter-lagann-announce.service

fleetctl start ./run/units/flitter-builder.service
fleetctl start ./run/units/flitter-builder-announce.service

fleetctl start ./run/units/flitter-havok.service

# Routing
fleetctl start ./run/units/flitter-router@1.service
