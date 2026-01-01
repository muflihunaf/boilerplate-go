# syntax=docker/dockerfile:1

# =============================================================================
# Build stage: Compile Go application
# =============================================================================
FROM golang:1.22-alpine AS builder

# Build arguments
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o /build/api \
    ./cmd/api

# =============================================================================
# Test stage: Run tests (optional, use with --target=test)
# =============================================================================
FROM builder AS test

RUN go test -v -race ./...

# =============================================================================
# Final stage: Minimal production image
# =============================================================================
FROM alpine:3.20 AS production

# Labels for container metadata
LABEL maintainer="muflihunaf" \
      org.opencontainers.image.title="Boilerplate Go API" \
      org.opencontainers.image.description="Production-ready Go API boilerplate" \
      org.opencontainers.image.source="https://github.com/muflihunaf/boilerplate-go"

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /build/api /app/api

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser:appgroup

# Environment variables (can be overridden at runtime)
ENV APP_ENV=production \
    PORT=8080 \
    TZ=UTC

# Expose application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["/app/api"]

# =============================================================================
# Distroless variant: Even smaller image (use with --target=distroless)
# =============================================================================
FROM gcr.io/distroless/static-debian12:nonroot AS distroless

LABEL maintainer="muflihunaf" \
      org.opencontainers.image.title="Boilerplate Go API" \
      org.opencontainers.image.description="Production-ready Go API boilerplate (distroless)"

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/api /app/api

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Environment variables
ENV APP_ENV=production \
    PORT=8080

# Expose application port
EXPOSE 8080

# Run as non-root (distroless nonroot user)
USER nonroot:nonroot

ENTRYPOINT ["/app/api"]
