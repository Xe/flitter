#!/bin/bash

set -ex

for file in ./drone/*
do
	echo "entering $file"
	if ! bash $file
	then
		exit 1
	fi
done

killall5 -9
