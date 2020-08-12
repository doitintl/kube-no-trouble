FROM golang:1.13.9-alpine3.11 as builder

WORKDIR /src

COPY go.mod .
COPY go.sum .
COPY scripts scripts

RUN scripts/alpine-setup.sh 

RUN go mod download

COPY . .

RUN make all

FROM scratch

COPY --from=builder /src/bin/kubent-linux-amd64 /app/kubent

WORKDIR /app

CMD ["./kubent"]
