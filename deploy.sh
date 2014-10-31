#!/bin/sh

set -ex

# deploy flitter
fleetctl start ./run/units/flitter-registry.service
fleetctl start ./run/units/flitter-announce@registry.service

fleetctl start ./run/units/flitter-lagann.service
fleetctl start ./run/units/flitter-announce@lagann.service

fleetctl start ./run/units/flitter-builder.service
fleetctl start ./run/units/flitter-announce@builder.service

fleetctl start ./run/units/flitter-havok.service

# Routing
fleetctl start ./run/units/flitter-router@1.service

# Proxy
fleetctl start ./run/units/flitter-proxy.service
