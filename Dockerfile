FROM golang:1.19.2-alpine3.16 as builder
ARG GITHUB_REF GITHUB_SHA
WORKDIR /src
COPY go.mod go.sum ./
COPY scripts scripts
RUN scripts/alpine-setup.sh
RUN go mod download
COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile
RUN export
RUN make all

FROM scratch
USER 10000:10000
WORKDIR /app
COPY --from=builder /src/bin/kubent-linux-amd64 /app/kubent
ENTRYPOINT ["/app/kubent"]
