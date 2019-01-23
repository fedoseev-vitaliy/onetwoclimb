# onetwoclimb
12Climb backend server

Prerequisites:
1. Install go `https://golang.org/doc/install`
2. Install make `https://brewformulas.org/Make` for MacOS, Ubuntu `sudo apt-get install make`
3. Clone repository to GOPATH/src/github.com/${onetwoclimb}

4. Compile binary with executing make cmd from repository root folder `make`
5. Run server from repository root folder `./onetwoclimb server --host 127.0.0.1 --port 8000`

To generate swagger doc run `make swaggerdoc`