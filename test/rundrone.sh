#!/bin/bash

set -ex

for file in ./drone/*
do
	bash $file
done
