FROM golang

WORKDIR /client

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY *.json ./

RUN go build
