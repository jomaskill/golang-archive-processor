FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0

RUN go build -o main .
