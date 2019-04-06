FROM golang:1.12.1

RUN apt-get update && apt-get install make

WORKDIR $GOPATH/src/github.com/onetwoclimb

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./vendor ./vendor
COPY Makefile ./main.go ./

RUN make build && \
    cp ./onetwoclimb /usr/local/bin/ && \
    rm -rf /go/src/github.com

WORKDIR /usr/local/bin/

ENV BIND 0.0.0.0:80

EXPOSE 80

ENTRYPOINT ["onetwoclimb"]

