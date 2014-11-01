# 0.0.1

Initial release, very immature.

 - Accepts repos containing a Dockerfile over `git push` then builds them into
   an image, spinning it up with direct docker API calls to the local host.
 - Announces new app and backend containers to etcd for vulcand to pick up on
 - TCP proxy for the builder ssh

# 0.0.2

 - fleet support instead of just shooting things to docker on Lagann's host
 - Lagann does not handle deployment anymore
