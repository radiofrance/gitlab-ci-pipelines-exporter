FROM docker.io/golang:1.19.4 as builder
ARG PLATFORM=linux
ARG ARCH=amd64
ARG VERSION=devel

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOOS=${PLATFORM} GOARCH=${ARCH} go build -ldflags="-w -s -X 'main.version=${VERSION}'" ./cmd/gitlab-ci-pipelines-exporter

FROM scratch

COPY --from=builder /build/gitlab-ci-pipelines-exporter /gitlab-ci-pipelines-exporter

EXPOSE 9252
USER 65534

ENTRYPOINT ["/gitlab-ci-pipelines-exporter"]
