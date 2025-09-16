# ---------- Stage 1: Build ----------
FROM golang:1.24.4-alpine AS builder

LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      stage="builder"

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
      -ldflags="-s -w" \
      -o /app/pcc-server \
      ./cmd/pcc-server

# ---------- Stage 2: Runtime ----------
FROM alpine:3.18

LABEL maintainer="Luka Usljebrka <lukauslje13@gmail.com>" \
      app="payment-card-center-service" \
      description="Payment Card Center service"

RUN apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates || true

RUN addgroup -S app && adduser -S -G app -u 10001 app

WORKDIR /app

# Kopiranje binarnog fajla i konfiguracije
COPY --from=builder /app/pcc-server /app/pcc-server
COPY --from=builder /src/config.yaml /app/config.yaml

# Kopiranje TLS sertifikata direktno u image
COPY certs/certs/pcc.crt /etc/ssl/certs/pcc.crt
COPY certs/private/pcc.key /etc/ssl/private/pcc.key

# Prava pristupa
RUN chown -R app:app /app /etc/ssl

USER app

ENV APP_ENV=production \
    PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/pcc-server"]