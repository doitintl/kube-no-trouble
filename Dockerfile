FROM --platform=$BUILDPLATFORM golang:1.22-alpine3.19 AS builder
ARG GITHUB_REF GITHUB_SHA
WORKDIR /src
COPY go.mod go.sum ./
COPY scripts scripts
RUN scripts/alpine-setup.sh
RUN go mod download
COPY cmd cmd
COPY pkg pkg
COPY Makefile Makefile

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make all

FROM scratch AS kubent
USER 10000:10000
WORKDIR /app
ARG TARGETOS TARGETARCH
COPY --from=builder /src/bin/kubent-$TARGETOS-$TARGETARCH /app/kubent
ENTRYPOINT ["/app/kubent"]
