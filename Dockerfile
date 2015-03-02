FROM golang:1.4

COPY . /go/src/github.com/nickschuch/marco
WORKDIR /go/src/github.com/nickschuch/marco

RUN go build

EXPOSE 80
ENTRYPOINT ["marco"]
CMD ["--help"]
