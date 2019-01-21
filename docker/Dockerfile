FROM golang:1.8

WORKDIR /go/src/github.com/cglotr/demos-golang/docker
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["docker"]
