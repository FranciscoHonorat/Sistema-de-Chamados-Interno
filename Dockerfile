# syntax=docker/dockerfile:1.7
#
# One image for the whole app: the Go binary serves the API and the built Vue
# frontend. Both build stages run on the build machine's platform and the Go
# stage cross-compiles, so multi-arch builds (amd64 + arm64) need no emulation.

ARG GO_VERSION=1.26
ARG NODE_VERSION=24

FROM --platform=$BUILDPLATFORM node:${NODE_VERSION}-alpine AS frontend
WORKDIR /src
COPY frontend/package.json frontend/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci --no-audit --no-fund
COPY frontend/ ./
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS backend
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=unknown
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY backend/ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o /out/sys-called ./cmd/sys-called

# distroless: no shell, no package manager, runs as an unprivileged user.
FROM gcr.io/distroless/static-debian12:nonroot
ARG VERSION=dev
ARG COMMIT=unknown
LABEL org.opencontainers.image.title="codeticket" \
      org.opencontainers.image.description="Sistema de controle de chamados internos (monólito modular Go + Vue)" \
      org.opencontainers.image.source="https://github.com/FranciscoHonorat/Sys-Called" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}"

COPY --from=backend /out/sys-called /usr/local/bin/sys-called
COPY --from=frontend /src/dist /app/web

# Secure by default: production mode requires a JWT key and never seeds the
# demo users. docker-compose.yml switches to development explicitly.
ENV APP_ENV=production \
    PORT=8080 \
    METRICS_PORT=9090 \
    STATIC_DIR=/app/web

USER nonroot:nonroot
EXPOSE 8080 9090

HEALTHCHECK --interval=10s --timeout=3s --start-period=20s --retries=5 \
    CMD ["/usr/local/bin/sys-called", "healthcheck"]

ENTRYPOINT ["/usr/local/bin/sys-called"]
CMD ["serve"]
