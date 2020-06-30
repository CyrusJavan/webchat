FROM golang:1.14

WORKDIR /go/src/app
COPY ./app .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["webchat"]