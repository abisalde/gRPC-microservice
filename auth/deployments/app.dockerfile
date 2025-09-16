# Build stage

FROM golang:1.25-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Install git for go modules if private repos
RUN apk add --no-cache git gcc musl-dev bash

RUN apt-get update && apt-get install -y curl unzip \
    && curl -L https://github.com/fullstorydev/grpcurl/releases/download/v1.9.1/grpcurl_1.9.1_linux_x86_64.tar.gz \
    | tar -xz -C /usr/local/bin

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


# Copy the entire auth service code
COPY auth/ ./auth/

# Copy air config from air.toml
COPY auth/deployments/air.toml ./auth/.air.toml

# Set working directory to auth service
WORKDIR /app/auth

ENV AIR_WORKSPACE=/app/auth \
    ENVIRONMENT=development



# Create directory for air builds
RUN mkdir -p /app/auth/tmp

# Expose the port your server runs on
EXPOSE 50051


CMD ["air", "-c", "/app/auth/.air.toml"]