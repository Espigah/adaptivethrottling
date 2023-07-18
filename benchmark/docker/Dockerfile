# syntax=docker/dockerfile:1

FROM golang:1.20

ENV GOOS linux
ENV GOARCH ${GOARCH:-amd64}
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./


RUN go build -v -o api cmd/api/main.go

CMD ["api]