# use ubuntu focal as base image
# builder stage
FROM ubuntu:focal AS builder

# make sure we're root
USER root

# get build dependencies
# get go toolchain
WORKDIR /tmp
RUN apt-get update && apt-get install wget unzip -y
# we're not using ADD since it disables caching completely
RUN wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz -O /tmp/go.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go.linux-amd64.tar.gz
RUN rm /tmp/go.linux-amd64.tar.gz

# download purecss
WORKDIR /home/root/go/src/github.com/david-sorm/goblog/
COPY scripts/purecssInstall.sh purecssInstall.sh
RUN chmod +x purecssInstall.sh
RUN ./purecssInstall.sh

# copy source files
COPY . .

# get dependencies and compile
RUN /usr/local/go/bin/go get -d -v ./...
RUN /usr/local/go/bin/go build -o serve .
RUN chmod +x /home/root/go/src/github.com/david-sorm/goblog/run

# final image stage
FROM ubuntu:focal

# copy artefacts and needed files
RUN mkdir /app && mkdir /app/html
COPY --from=builder /home/root/go/src/github.com/david-sorm/goblog/serve /app/serve
COPY --from=builder /home/root/go/src/github.com/david-sorm/goblog/html/ /app/html/
COPY scripts/.docker-conf-gen.sh /app/docker-conf-gen.sh
COPY scripts/.docker-run.sh /app/docker-run.sh

# open port
EXPOSE 80

# register all args
ENV BLOGNAME="" ARTICLESPERPAGE=5 LISTENON=80 STORE="mock" STORE_HOST="" STORE_DB="" STORE_USER="" STORE_PASSWORD="" CACHINGENGINE="" HOTSWAPTEMPLATES="no"

# run
WORKDIR /app
RUN chmod +x docker-run.sh
RUN chmod +x docker-conf-gen.sh
CMD ./docker-run.sh
