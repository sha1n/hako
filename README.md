[![Build Status](https://travis-ci.org/sha1n/hako.svg?branch=master)](https://travis-ci.org/sha1n/hako) [![Go Report Card](https://goreportcard.com/badge/sha1n/hako)](https://goreportcard.com/report/sha1n/hako)

# hako
Just a simple HTTP echo server with CLI interface for testing...


## Setting up from sources
```bash
git clone git@github.com:sha1n/hako.git
cd hako

# build the app
make

# optionally copy the binary to your path
cp bin/hako <~/.local/bin>
```

## Downloading released binaries

**MacOS cURL Example**
```bash
curl -Lf --compressed -o ~/.local/bin/hako https://github.com/sha1n/hako/releases/download/v0.3.0/hako-darwin-amd64
```

## Usage Example
Terminal A:
```bash 
# run the server
hako start -p 80 --verbose --path /echo/shmecho
2020/03/15 15:38:45 Registering signal listeners for graceful HTTP server shutdown..
2020/03/15 15:38:45 Staring HTTP Server on :80
2020/03/15 15:38:45 Waiting for shutdown signal...
2020/03/15 15:39:02 Handling request at /echo/shmecho
Received: {'Hello': 'World'}
```
Terminal B:
```bash 
curl -X POST localhost:80/echo/shmecho -H "Content-Type: application/json" --data "{'Hello': 'World'}"
```
