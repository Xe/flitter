flitter
=======

[![GoDoc](https://godoc.org/github.com/Xe/flitter?status.svg)](https://godoc.org/github.com/Xe/flitter) [![Build Status](https://drone.io/github.com/Xe/flitter/status.png)](https://drone.io/github.com/Xe/flitter/latest)

A minimal platform-as-a-service using CoreOS, fleet, vulcand, and docker-havok.

Flitter is made up of many parts from many different authors. Where possible 
all existing code is kept under the terms of the license it came from. Any new 
projects inside this repository are under the highly permissive Zlib license:

```
Copyright (C) 2014 Sam Dodrill <xena@yolo-swag.com> All rights reserved.

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgment in the product documentation would be
   appreciated but is not required.

2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.

3. This notice may not be removed or altered from any source
   distribution.
```

If you find a program is incorrectly licensed please open a github issue so it 
can be fixed as soon as possible. These kinds of issues are critical and will 
be treated as such.

## Installing

Use `./deploy.sh`. Future tooling to automate the editing of the units for 
custom deployments will be present in a future release.

### Havok

You must set a domain on the flitter host using:

```console
$ etcdctl set /flitter/domain flitter.xeserv.us
```

If you are using Vagrant (`contrib/vagrant`) to test Flitter, you can use
`flitter.xeserv.us` as your domain setting. Everything will be forwarded 
properly.

### Port Forwarding

Allow ports `80`, `22`, and `2232` from any IP address.

## Support

Flitter is **PRE-ALPHA** software. It may eat your hamster. If you use this in 
production as is, the authors take **NO** fault whatsoever.

At this stage Flitter is known to be able to do `git push` deploys to a cluster 
of servers running [fleet](http://github.com/coreos/fleet).

[![](https://i.imgur.com/Et0muSJ.png)](https://asciinema.org/a/13457)

---

A public project by XeServ.
