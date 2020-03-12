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
