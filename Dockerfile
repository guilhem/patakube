FROM golang:1.7

ADD . /go/src/github.com/guilhem/patakube

RUN go install github.com/guilhem/patakube

EXPOSE 8080

ENTRYPOINT ["/go/bin/patakube"]
