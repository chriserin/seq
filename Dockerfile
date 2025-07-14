# Multi-stage Dockerfile for static compilation of seq with Lua and RtMidi
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    g++ \
    musl-dev \
    alsa-lib-dev \
    lua5.4-dev \
    pkgconfig \
    make

# Create symlink for lua5.4 static library so the linker can find it
RUN ln -s /usr/lib/lua5.4/liblua.a /usr/lib/liblua5.4.a

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Copy vendor directory if it exists
COPY vndr/ ./vndr/

# Download dependencies
RUN go mod download

# Copy Go source files
COPY *.go ./
COPY internal/ ./internal/

# Set environment variables for mixed linking (static Lua, dynamic ALSA)
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-L/usr/lib/lua5.4 -Wl,-Bstatic -llua5.4 -Wl,-Bdynamic -lm"

# Build the application
RUN go build -tags lua54 \
    -ldflags "-s -w" \
    -o seq

# Verify linking
RUN ldd seq

# Final stage - runtime image with ALSA only (Lua is static)
FROM alpine:latest

# Install runtime dependencies (ALSA + C++ runtime for RtMidi, Lua is static)
RUN apk add --no-cache \
    alsa-lib \
    libstdc++ \
    libgcc

# Copy the binary
COPY --from=builder /app/seq /seq

# Set the binary as the entrypoint
ENTRYPOINT ["/seq"]
