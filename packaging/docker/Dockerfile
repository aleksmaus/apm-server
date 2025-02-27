ARG BASE_IMAGE=ubuntu:20.04

# Build the apm-server binary. The golang image version is kept
# up to date with .go-version by Makefile.
FROM golang:1.18.8 as builder
WORKDIR /src
COPY go.mod go.sum /src/
COPY internal/approvaltest/go.mod /src/internal/approvaltest/go.mod
COPY internal/glog/go.mod /src/internal/glog/go.mod
COPY internal/otel_collector/go.mod /src/internal/otel_collector/go.mod
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY Makefile *.mk /src
COPY cmd /src/cmd
COPY internal /src/internal
COPY x-pack /src/x-pack
COPY .git /src/.git

ENV GOROOT /usr/local/go
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make apm-server

FROM ${BASE_IMAGE}
ARG TARGETARCH
ARG BUILD_DATE
ARG VERSION
ARG VCS_REF

# Add an init process.
ADD --chmod=0755 packaging/docker/tini-${TARGETARCH}.tar.xz /usr/bin

# Statically defined labels.
LABEL \
  org.label-schema.schema-version="1.0" \
  org.label-schema.vendor="Elastic" \
  org.label-schema.license="Elastic License" \
  org.label-schema.name="apm-server" \
  org.label-schema.url="https://www.elastic.co/apm" \
  org.label-schema.vcs-url="github.com/elastic/apm-server" \
  io.k8s.description="Elastic APM Server" \
  io.k8s.display-name="Apm-Server image" \
  org.opencontainers.image.licenses="Elastic License" \
  org.opencontainers.image.title="Apm-Server" \
  org.opencontainers.image.vendor="Elastic" \
  name="apm-server" \
  maintainer="infra@elastic.co" \
  vendor="Elastic" \
  release="1" \
  url="https://www.elastic.co/apm" \
  summary="apm-server" \
  license="Elastic License" \
  description="Elastic APM Server"

# Dynamic labels, only set in published images.
LABEL \
  org.label-schema.build-date=${BUILD_DATE} \
  org.label-schema.version=${VERSION} \
  org.label-schema.vcs-ref=${VCS_REF} \
  org.opencontainers.image.created=${BUILD_DATE} \
  version=${VERSION}

ENV ELASTIC_CONTAINER "true"
ENV PATH=/usr/share/apm-server:$PATH

# When running under Docker, we must ensure libbeat monitoring pulls cgroup
# metrics from /sys/fs/cgroup/<subsystem>/, ignoring any paths found in
# /proc/self/cgroup.
ENV LIBBEAT_MONITORING_CGROUPS_HIERARCHY_OVERRIDE=/

COPY --chmod=0755 packaging/docker/docker-entrypoint /usr/local/bin/docker-entrypoint
COPY --chmod=0644 licenses/ELASTIC-LICENSE-2.0.txt NOTICE.txt /licenses/

# NOTE(axw) apm-server.yml and data/ must be owned by apm-server
# to allow the image to be run with a non-root user.
WORKDIR /usr/share/apm-server
COPY --chmod=0640 --chown=1000:0 apm-server.yml ./apm-server.yml
COPY --chmod=0755 --from=builder /src/apm-server ./apm-server
RUN sed -i 's/localhost:8200/0.0.0.0:8200/' apm-server.yml
RUN sed -i 's/localhost:9200/elasticsearch:9200/' apm-server.yml
RUN mkdir --mode=0775 data && chown 1000 data
RUN echo 'apm-server:*:1000:0::/usr/share/apm-server:/bin/false' >> /etc/passwd

USER apm-server
EXPOSE 8200
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/docker-entrypoint"]
