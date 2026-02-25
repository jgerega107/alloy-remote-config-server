# Build stage
FROM golang:1.26-alpine AS builder

# Build arguments for multi-arch support
ARG TARGETOS
ARG TARGETARCH

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for build stage (for consistency)
RUN adduser -D -g '' appuser

# Set working directory with correct module path
WORKDIR /src

# Copy all source code
COPY . .

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod tidy && go mod download

# Build the application with optimizations
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -ldflags="-w -s" -trimpath \
    -o /alloy-remote-config-server cmd/config/main.go

# Final stage - minimal image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' -h /app -s /bin/false appuser && \
    mkdir -p /app/conf && \
    chown -R appuser:appuser /app

# Copy binary from builder
COPY --from=builder /alloy-remote-config-server /usr/local/bin/alloy-remote-config-server

# Set ownership
RUN chown appuser:appuser /usr/local/bin/alloy-remote-config-server

# Switch to non-root user
USER appuser

# Set working directory
WORKDIR /app

# Expose ports (documentation purposes)
EXPOSE 8080 8888

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["/usr/local/bin/alloy-remote-config-server"]
