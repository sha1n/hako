[![Build Status](https://travis-ci.org/sha1n/echo-server.svg?branch=master)](https://travis-ci.org/sha1n/echo-server) [![Go Report Card](https://goreportcard.com/badge/sha1n/echo-server)](https://goreportcard.com/report/sha1n/echo-server)

# echo-server
Just a simple HTTP echo server with CLI interface for testing...


## Setting up from sources
```bash
git clone git@github.com:sha1n/echo-server.git
cd echo-server

# build the app
make

# optionally copy the binary to your path
cp bin/echoserver <~/.local/bin>

# run the server
echoserver start -p 80 --path /echo/shmecho
```

## Downloading released binaries

**MacOS cURL Example**
```bash
curl -Lf --compressed -o ./echoserver https://github.com/sha1n/echo-server/releases/download/v0.1/echo-server-darwin-amd64
```
