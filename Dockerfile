FROM docker.io/golang:1.22.1 as builder
ARG VERSION=devel

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.version=${VERSION}'" ./cmd/gitlab-ci-pipelines-exporter

FROM scratch

COPY --from=builder /build/gitlab-ci-pipelines-exporter /gitlab-ci-pipelines-exporter

EXPOSE 9252
USER 65534

ENTRYPOINT ["/gitlab-ci-pipelines-exporter"]
