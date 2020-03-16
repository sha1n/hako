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
cp bin/hako <~/.local/bin/hako>
```

## Downloading released binaries

**MacOS cURL Example**
```bash
curl -Lf --compressed -o ~/.local/bin/hako https://github.com/sha1n/hako/releases/download/v0.3.0/hako-darwin-amd64

chmod +x ~/.local/bin/hako
```

## Usage Example

**Terminal A:**
```bash 
# run the server
➜  ~ hako start --verbose --port 8090 --path /echo/shmecho
[HAKO] 2020/03/16 09:46:34 Registering signal listeners for graceful HTTP server shutdown..
[HAKO] 2020/03/16 09:46:34 Staring HTTP Server on :8090
[HAKO] 2020/03/16 09:46:34 Waiting for shutdown signal...
[HAKO] 2020/03/16 09:53:39 Handling request at /echo/shmecho
[HAKO] 2020/03/16 09:53:39 Body: {'Hello': 'World'}

[GIN] 2020/03/16 - 09:53:39 | 200 |      54.324µs |             ::1 | POST     /echo/shmecho
[HAKO] 2020/03/16 09:53:41 Handling request at /non-existing
[GIN] 2020/03/16 - 09:53:41 | 404 |      12.785µs |             ::1 | HEAD     /non-existing
```

**Terminal B:**
```bash 
# posting to an existing URL
➜  ~ curl -X POST localhost:8090/echo/shmecho -H "Content-Type: application/json" --data "{'Hello': 'World'}"
{'Hello': 'World'}%

# heading to a non-existing URL
➜  ~ curl -I localhost:8090/non-existing
HTTP/1.1 404 Not Found
Content-Type: text/plain
Date: Mon, 16 Mar 2020 07:53:41 GMT
Content-Length: 18
```
