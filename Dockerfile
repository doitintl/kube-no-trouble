FROM golang:1.18.3-alpine3.16 as builder

WORKDIR /src

COPY go.mod .
COPY go.sum .
COPY scripts scripts

RUN scripts/alpine-setup.sh

RUN go mod download

COPY . .

RUN make all

FROM alpine:3.16.0

USER 1000

COPY --from=builder /src/bin/kubent-linux-amd64 /app/kubent

WORKDIR /app
