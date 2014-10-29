#!/bin/bash

set -ex

export KEY=$(cat ~/.ssh/id_rsa.pub)

export code=$(curl --data '{"name": "drone", "sshkeys": [{"comment": "drone", "key": "'"$KEY"'", "fingerprint": "b3:79:71:a8:f5:18:af:e3:da:d7:a4:5e:db:03:ac:80"}]}' -X POST http://127.0.0.1:3000/register --write-out "%{http_code}\n" --silent --output /dev/null)
if ! [ "$code" = "200" ]; then
	echo "Cannot create account :("
	exit 1
fi

export code=$(curl --data '{"name": "drone", "sshkeys": [{"comment": "drone", "key": "'"$KEY"'", "fingerprint": "b3:79:71:a8:f5:18:af:e3:da:d7:a4:5e:db:03:ac:80"}]}' -X POST http://127.0.0.1:3000/register --write-out "%{http_code}\n" --silent --output /dev/null)
if ! [ "$code" = "409" ]; then
	echo "Can create duplicate accounts"
	exit 1
fi
