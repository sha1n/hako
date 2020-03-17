[![Build Status](https://travis-ci.org/sha1n/hako.svg?branch=master)](https://travis-ci.org/sha1n/hako) [![Go Report Card](https://goreportcard.com/badge/sha1n/hako)](https://goreportcard.com/report/sha1n/hako) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/sha1n/hako.svg?style=flat-square)](https://github.com/sha1n/hako/releases)

# Hako
Hako is an HTTP echo server with a CLI interface that provides some extra features. Hako has been developed to help me develop and test one of my projects and since the need for such utility keeps popping every few years, I thought it would be nice to share it with others.


## Building from sources
```bash
git clone git@github.com:sha1n/hako.git
cd hako

# build the Go app (a local Go installation is required)
make

# optionally copy the binary to your path
cp bin/hako <~/.local/bin/hako>
```

## Downloading released binaries

**MacOS cURL Example**
```bash
curl -Lf --compressed -o <~/.local/bin/hako> https://github.com/sha1n/hako/releases/download/v0.5.0/hako-darwin-amd64

chmod +x <~/.local/bin/hako>
```

## Usage Example
See usage examples below. Use `hako --help` for help.

**Terminal A:**
```bash 
# run the server
➜  ~ hako start -p 8090 --path /echo/shmecho --delay 1 --verbose --verbose-headers
[HAKO] 2020/03/17 12:32:36 Registering signal listeners for graceful HTTP server shutdown..
[HAKO] 2020/03/17 12:32:36 Staring HTTP Server on :8090
[HAKO] 2020/03/17 12:32:36 Waiting for shutdown signal...
[HAKO] 2020/03/17 12:32:38 Handling request at /echo/shmecho
[HAKO] 2020/03/17 12:32:38 Received headers:

User-Agent : curl/7.64.1
Accept : */*
Content-Type : application/json
Content-Length : 18

[HAKO] 2020/03/17 12:32:38 Received body:

{'Hello': 'World'}

[HAKO] 2020/03/17 12:32:38 Delaying response in 1 millis
[GIN] 2020/03/17 - 12:32:38 | 200 |    1.328925ms |             ::1 | POST     /echo/shmecho
[HAKO] 2020/03/17 12:33:44 Handling request at /non-existing
[GIN] 2020/03/17 - 12:33:44 | 404 |      14.822µs |             ::1 | HEAD     /non-existing
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
Date: Tue, 17 Mar 2020 10:33:44 GMT
Content-Length: 18
```
