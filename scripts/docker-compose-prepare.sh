#!/bin/bash
# just creates folders needed for the example docker-example stack

# make sure we're not in the 'scripts' directory
if [[ $(pwd) == */scripts ]]
then
  cd ..
fi

mkdir -p docker-data/postgres