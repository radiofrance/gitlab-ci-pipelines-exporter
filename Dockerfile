FROM docker.io/golang:1.27.3@sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e AS builder
ARG VERSION=devel

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.version=${VERSION}'" ./cmd/gitlab-ci-pipelines-exporter

FROM scratch

COPY --from=builder /build/gitlab-ci-pipelines-exporter /gitlab-ci-pipelines-exporter

EXPOSE 9252
USER 65534

ENTRYPOINT ["/gitlab-ci-pipelines-exporter"]
