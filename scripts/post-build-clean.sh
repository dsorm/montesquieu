#!/bin/bash
# Removes everything from source that's not needed after compilation

if [[ $1 != "yes" ]]
then
  echo "Do you want to remove source files that are not needed after compilation? If you're sure, run ./post-build-clean.sh yes"
  exit
fi

# make sure we're not in the 'scripts' directory
if [[ $(pwd) == */scripts ]]
then
  cd ..
fi

# since we like to live dangerously
rm -rf .git .idea .github article config globals handlers run template .gitignore run.go go.mod readme.md .env go.sum docker.config.json | true
