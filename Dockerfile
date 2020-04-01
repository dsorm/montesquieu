# use golang with debian buster as base image
FROM golang:buster

# copy source files
COPY . /opt/goblog/.

# delete useless crap
RUN rm /opt/goblog/readme.md || true

# delete config.json, if exists, and copy over the docker config file
RUN rm /opt/goblog/config.json || true
COPY docker.config.json /opt/goblog/config.json

# get dependencies and compile
WORKDIR /opt/goblog
RUN go get -d -v ./...
RUN go build -o run .

# open port
EXPOSE 80

# run
CMD ["/opt/goblog/run"]