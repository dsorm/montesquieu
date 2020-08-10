# project goblog (working name)
a CS50 final project, and yeah, I know, the name's really bad

it doesn't really have any functionality at the moment, it uses a mockup backend now

## How to use
### With Docker & Docker Compose (the recommended way)
- Use a Linux system (something like Ubuntu/Debian, CentOS/Fedora/RHEL will work) with Docker and Docker Compose installed, there is no official support for other platforms
- Download a copy of master to your system, open a terminal and navigate to the downloaded folder
- `docker build -t goblog .`
- If the command above fails, run it with sudo, so `sudo docker build -t goblog .`
- You might want to take a look at the docker-compose.yml file and change some environment variables (hint: look at the comments) before continuing
- `docker-compose up -d`
- Goblog should be running on localhost:80, unless you changed the port
- If you want to stop goblog, just run `docker-compose down` from the same folder
### Without Docker on Linux or similar: (less recommended)
- Make sure you have a recent version of go toolchain installed
- Clone master / download a [release] of goblog
- Install dependencies using `go get ./..` within the directory
- Build the executable using `go build -o run .`
- Run the executable: `./run`
- Config.json with default settings will be made on first startup, you can change any of the settings and restart
### Without Docker on Windows: (least recommended)
- Same as on Linux, just instead of `go build -o run .` use  `go build -o run.exe` and start `run.exe` instead of doing `./run`
### With Docker:
- It's recomended to run everything as sudo/root, since only root can map ports below 1024
- Clone master / download a [release] of goblog
- Change the docker.config.json according to your needs first (except for ListenOn)
- Run this command within the directory to create docker image: `docker build -t goblog .`
- Afterward, change `<your_port>` to your own preffered port and run this command to start a container with the goblog image we created earlier: `docker run -d -p 8080:<your_port> --name goblog localhost/goblog` 

[release]: https://github.com/david-sorm/goblog/releases