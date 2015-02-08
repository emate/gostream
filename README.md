# gostream
Simple server for streaming local files to remote clients via TCP. 


### Installation:

gostream depends on github.com/ActiveState/tail module

```sh
$ go get https://github.com/ActiveState/tail
```

Get gostream and build it!
```sh
$ git clone https://github.com/emate/gostream 
$ go build gostream.go
```

### How to run:
Run server
```sh
$ gostream [-l 0.0.0.0:8080] <filename>
```

On the other side, you can connect to socket and receive streamed file:
```sh
$ telnet REMOVE_GOSTREAM_SERVER_IP PORT
```

