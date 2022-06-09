FROM golang

WORKDIR /client

COPY wait-for-it.sh ./wait-for-it.sh
RUN chmod +x ./wait-for-it.sh

COPY wait-for-service-readiness.sh ./wait-for-service-readiness.sh
RUN chmod +x ./wait-for-service-readiness.sh

WORKDIR /client

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.json ./

RUN go build
#ENTRYPOINT ["/bin/sh","-c","./wait-for-it.sh accountapi:8080 -t 30"]
#CMD ["go", "test"]