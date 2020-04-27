# use ubuntu focal as base image
FROM ubuntu:focal

# make sure we're root
USER root

# download apt package info
RUN apt-get update

# get build dependencies
# get go toolchain
WORKDIR /tmp
RUN apt-get install wget unzip -y
RUN wget "https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz
RUN rm /tmp/go1.14.2.linux-amd64.tar.gz

# download purecss
WORKDIR /home/root/go/src/github.com/david-sorm/goblog/
COPY purecssInstall.sh purecssInstall.sh
RUN chmod +x purecssInstall.sh
RUN ./purecssInstall.sh

# copy source files
COPY . .

# delete useless crap
RUN rm readme.md || true
RUN rm config.json || true

# get dependencies and compile
RUN /usr/local/go/bin/go get -d -v ./...
RUN /usr/local/go/bin/go build -o serve .
RUN chmod +x /home/root/go/src/github.com/david-sorm/goblog/run

# delete build dependencies
RUN rm -r /usr/local/go/*
RUN rmdir /usr/local/go
RUN apt-get autoremove -y

# open port
EXPOSE 8080

# run
CMD su root -c /home/root/go/src/github.com/david-sorm/goblog/serve