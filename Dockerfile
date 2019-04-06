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

# TODO make multistage build
#FROM alpine:latest
#WORKDIR /root/
#RUN apk --no-cache add ca-certificates && mkdir /static
#COPY --from=builder /usr/local/bin .
#RUN pwd && ls -v

ENV BIND 0.0.0.0:80
EXPOSE 80
ENTRYPOINT ["onetwoclimb"]

