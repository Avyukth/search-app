# Build Stage
FROM  golang:1.21-bullseye as builder

# Set Environment Variables
ENV CGO_ENABLED=0
ENV GOCACHE=/root/.cache/go-build
ARG BUILD_REF

# Set Work Directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies and cache them
RUN --mount=type=cache,target=/root/.cache/go-build go mod tidy

# Copy the entire directory
COPY . .

# Build the application
WORKDIR /app/cmd/server
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags "-X main.build=${BUILD_REF}"

# Runtime Stage
FROM gcr.io/distroless/static-debian11

# Metadata as described in OCI image spec annotations
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
    org.opencontainers.image.title="search-app" \
    org.opencontainers.image.authors="Subhrajit Makur <makur.subhrajit@gmail.com>" \
    org.opencontainers.image.source="https://github.com/avyukth/search-app" \
    org.opencontainers.image.revision="${BUILD_REF}" \
    org.opencontainers.image.vendor="subhrajit.me Inc."

# Copy the binary from builder stage
COPY --from=builder /app/cmd/server/server /app/server
COPY --from=builder /app/.env /app/.env
# Set Work Directory
WORKDIR /app

# Command to run the application
CMD ["./server"]
