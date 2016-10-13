FROM golang:1.6
MAINTAINER "Quantum, Inc"
EXPOSE 8087

COPY . /go/src/github.com/quantum/discovery
RUN go install -v github.com/quantum/discovery

CMD ["discovery"]
