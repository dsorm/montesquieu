# project goblog (working name)
a CS50 final project, and yeah, I know, the name's really bad

## How to use
### Without Docker on Linux or similar:
- Make sure you have a recent version of go toolchain installed
- Clone master / download a [release] of goblog
- Install dependencies using `go get ./..` within the directory
- Build the executable using `go build -o run .`
- Run the executable: `./run`
### Without Docker on Windows:
- Same as on Linux, just instead of `go build -o run .` use  `go build -o run.exe` and start `run.exe` instead of doing `./run`
### With Docker:
- Clone master / download a [release] of goblog
- Change the docker.config.json according to your needs first (except for ListenOn)
- Run this command within the directory to create docker image: `docker build -t goblog .`
- Afterward, change `<your_port>` to your own preffered port and run this command to start a container with the goblog image we created earlier: `docker run -d -p 80:<your_port> --name goblog localhost/goblog` 

[release]: https://github.com/david-sorm/goblog/releases