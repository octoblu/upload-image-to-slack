FROM golang
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/upload-image-to-slack
COPY . /go/src/github.com/octoblu/upload-image-to-slack

RUN env CGO_ENABLED=0 go build -a -ldflags '-s' .

CMD ["./upload-image-to-slack"]
