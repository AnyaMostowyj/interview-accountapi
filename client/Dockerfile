FROM golang

WORKDIR /client

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY ./test_data/*.json ./test_data/

RUN go build
