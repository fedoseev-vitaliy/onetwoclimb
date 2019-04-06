FROM golang:1.12.1 as builder

RUN apt-get update && apt-get install make

WORKDIR $GOPATH/src/github.com/onetwoclimb

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./vendor ./vendor
COPY Makefile ./main.go ./

RUN make build && \
    cp ./onetwoclimb /usr/local/bin/ && \
    rm -rf /go/src/github.com

#FROM alpine
#COPY --from=builder /usr/local/bin /usr/local/bin
#
#WORKDIR /usr/local/bin
#
#RUN ls

ENV BIND 0.0.0.0:80
EXPOSE 80
ENTRYPOINT ["onetwoclimb"]

