FROM golang:1.4

COPY . /go/src/github.com/nickschuch/marco
WORKDIR /go/src/github.com/nickschuch/marco

RUN go get github.com/Sirupsen/logrus
RUN go get gopkg.in/alecthomas/kingpin.v1
RUN go get github.com/samalba/dockerclient
RUN go get github.com/nickschuch/go-tutum/tutum
RUN go get github.com/daryl/cash

RUN go build

EXPOSE 80
ENTRYPOINT ["marco"]
CMD ["--help"]
