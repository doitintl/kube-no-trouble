FROM 1.17.1-alpine3.14 as builder

WORKDIR /src

COPY go.mod .
COPY go.sum .
COPY scripts scripts

RUN scripts/alpine-setup.sh

RUN go mod download

COPY . .

RUN make all

FROM scratch

USER 1000

COPY --from=builder /src/bin/kubent-linux-amd64 /app/kubent

WORKDIR /app

CMD ["./kubent"]
