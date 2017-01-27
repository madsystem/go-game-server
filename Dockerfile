FROM golang
MAINTAINER madsystem@gmail.com

RUN go get github.com/madsystem/go-game-server
ENTRYPOINT /go/bin/go-game-server

EXPOSE 4444
EXPOSE 4446
