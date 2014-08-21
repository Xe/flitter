# SSH Daemon

## Startup

- [ ] Load SSH server keys
- [ ] Allocate listener port
- [ ] Run forever

## On Connection

- [ ] perform SSH handshake
- [ ] get user name they are connecting as
- [ ] get environment variables
- [ ] check ssh key fingerprint against etcd, fail!
- [ ] check repo against permissions list for user, fail!
- [ ] kick off builder
