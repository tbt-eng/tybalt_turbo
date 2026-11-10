# Multi-stage Dockerfile for Tybalt PocketBase Application with Litestream

# Stage 1: Build the Svelte UI
FROM node:24-alpine AS ui-builder

WORKDIR /app

# Install jq for version processing
RUN apk add --no-cache jq git

# Copy version files and scripts
COPY version.json ./
COPY scripts/version.sh ./scripts/
RUN chmod +x scripts/version.sh

WORKDIR /app/ui

# Copy package files
COPY ui/package.json ui/package-lock.json ./

# Install dependencies
RUN npm ci

# Copy UI source code
COPY ui/ ./

# Set the PocketBase URL for the production build
# This will be baked into the static build since SvelteKit uses adapter-static
ARG PUBLIC_POCKETBASE_URL
ENV PUBLIC_POCKETBASE_URL=${PUBLIC_POCKETBASE_URL}

# Build the UI
RUN npm run build

# Stage 2: Build the Go application
FROM golang:1.26.10-alpine AS go-builder

# Install git, jq and other build dependencies
RUN apk add --no-cache git ca-certificates jq

WORKDIR /app

# Accept git info as build args (only needed for Go stage)
ARG GIT_COMMIT
ARG GIT_COMMIT_SHORT
ARG GIT_BRANCH

# Set them as environment variables for the version script
ENV GIT_COMMIT=${GIT_COMMIT}
ENV GIT_COMMIT_SHORT=${GIT_COMMIT_SHORT}
ENV GIT_BRANCH=${GIT_BRANCH}

# Copy version files and scripts
COPY version.json ./
COPY scripts/version.sh ./scripts/
RUN chmod +x scripts/version.sh

# Copy go mod files
COPY app/go.mod app/go.sum ./

# Download dependencies
RUN go mod download

# Copy Go source code
COPY app/ ./

# Generate version information for the UI build
RUN cd /app && ./scripts/version.sh go > constants/version.go || \
    echo "Warning: Failed to generate version constants"

# Build the Help chat corpus from this commit's user docs, design notes and
# code. The helpcorpus package embeds it from helpcorpus/data, which is outside
# pb_public. A failure stops the build. --require-rules fails when a source file
# changed and its rule file in helpcorpus/rules did not: write the rules with
# `go run ./cmd/helprules` and commit them before you push.
COPY ui/src /corpus-src/ui/src
COPY descriptions /corpus-src/descriptions
# The files in helpcorpus/collections say who may read and write each
# PocketBase collection. This fails when a migration changed a collection and
# the files were not written again: go run ./cmd/helpcollections
RUN go run ./cmd/helpcollections --check
RUN go run ./cmd/helpcorpus --repo /corpus-src --app-dir . --require-rules

# Copy built UI from previous stage
COPY --from=ui-builder /app/ui/build ./pb_public

# Build the Go application
RUN CGO_ENABLED=0 go build -o tybalt main.go

# Stage 3: Final runtime image
FROM alpine:latest

# Install dependencies
RUN apk add --no-cache ca-certificates sqlite tzdata

# Install litestream
RUN wget https://github.com/benbjohnson/litestream/releases/download/v0.5.10/litestream-0.5.10-linux-x86_64.tar.gz -O - | tar -xzf - -C /usr/local/bin

# Create app directory
WORKDIR /app

# Copy the built application and UI
COPY --from=go-builder /app/tybalt ./tybalt

# Copy the UI assets
COPY --from=go-builder /app/pb_public ./pb_public

# Copy litestream configuration
COPY litestream.yml /etc/litestream.yml

# Copy startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Create pb_data directory for database persistence
RUN mkdir -p /app/pb_data

# Expose port 8080
EXPOSE 8080

# Use the startup script
CMD ["/start.sh"] 
