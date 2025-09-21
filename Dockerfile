# ---------- Stage 1: Build ----------
FROM golang:1.24-bullseye AS builder

LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      stage="builder"

WORKDIR /src

# copy go modules first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# allow cross-compilation
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# build binary (with CGO for TLS if needed later)
RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /app/pcc-server ./cmd/pcc-server

# ---------- Stage 2: Runtime ----------
FROM debian:bookworm-slim

LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      app="payment-card-center-service" \
      description="Payment Card Center service"

# install only required runtime packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# create non-root user
RUN groupadd -r app && useradd -r -g app -u 10001 app

WORKDIR /app

# copy binary and config
COPY --from=builder /app/pcc-server /app/pcc-server
COPY --from=builder /src/config.yaml /app/config.yaml

# fix ownership
RUN chown -R app:app /app

USER app

ENV APP_ENV=production \
    PORT=8081

EXPOSE 8081

ENTRYPOINT ["/app/pcc-server"]
