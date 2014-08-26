# SSH Daemon

## Startup

- [X] Load SSH server keys
- [X] Allocate listener port
- [X] Run forever

## On Connection

- [X] perform SSH handshake
- [X] get user name they are connecting as
- [X] get environment variables
- [X] check ssh key fingerprint against etcd, fail!
- [ ] check repo against permissions list for user, fail!
- [ ] kick off builder
