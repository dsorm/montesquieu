#!/bin/bash
# builds docker image from source

# make sure we're not in the 'scripts' directory
if [[ $(pwd) == */scripts ]]
then
  cd ..
fi

docker build -t montesquieu .