FROM golang

WORKDIR /client

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY ./test_data/*.json ./test_data/

COPY ./scripts/api_health.sh /
RUN chmod +x /api_health.sh

RUN go build
ENTRYPOINT ["/api_health.sh"]