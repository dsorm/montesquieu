#!/bin/bash
# used as the final CMD command for Dockerfile

# make a new onfig from environment variables
rm config.json | true
./docker-conf-gen.sh

# run
# `serve` is the artefact made from compiling montesquieu
./serve