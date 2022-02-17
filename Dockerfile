FROM golang:1.17-alpine AS staging

# Download goose migrations tool
# RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# 1. Precompile the entire go standard library into the first Docker cache layer: useful for other projects too!
RUN CGO_ENABLED=0 GOOS=linux go install -v -installsuffix cgo -a std

RUN apk update && apk add bash

WORKDIR $GOPATH/src/github.com/orlandoromo97/mailx-google-service

COPY go.mod go.sum /
RUN go mod download -x

COPY . .

# Building mailx-google-service binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo -o . ./cmd/mailx-google-service

COPY ./scripts/dev.sh /
RUN  chmod +x /dev.sh
ENTRYPOINT [ "/dev.sh" ]


