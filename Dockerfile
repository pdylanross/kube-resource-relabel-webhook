FROM golang:1.21 as build

ARG VERSION="0.0.0-docker"
ARG COMMIT_HASH="n/a-docker"
ARG BUILD_TIMESTAMP="n/a-docker"

ENV GO111MODULE=on

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/kube-resource-relabel-webhook \
    -ldflags="-X 'main.Version=$VERSION' -X 'main.CommitHash=$COMMIT_HASH' -X 'main.BuildTimestamp=$BUILD_TIMESTAMP'" \
    ./cmd/kube-resource-relabel

FROM gcr.io/distroless/base-debian12:latest

ARG VERSION="0.0.0-docker"
ARG COMMIT_HASH="n/a-docker"
ARG BUILD_TIMESTAMP="n/a-docker"

LABEL org.opencontainers.image.source="https://github.com/pdylanross/kube-resource-relabel-webhook"
LABEL org.opencontainer.image.description="A lightweight kubernetes resource relabeling mutation webhook"
LABEL org.opencontainers.image.created="$BUILD_TIMESTAMP"
LABEL org.opencontainers.image.version="$VERSION"
LABEL org.opencontainers.image.revision="$COMMIT_HASH"

COPY --from=build /go/bin/kube-resource-relabel-webhook /
ENTRYPOINT ["/kube-resource-relabel-webhook"]