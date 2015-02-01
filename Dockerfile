FROM golang:1.3

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# Get the libraries.
RUN go get github.com/Sirupsen/logrus
RUN go get gopkg.in/alecthomas/kingpin.v1
RUN go get github.com/samalba/dockerclient
RUN go get github.com/nickschuch/go-tutum/tutum

# Build the binary.
RUN go build

EXPOSE 80
ENTRYPOINT ["marco"]
CMD ["--help"]
