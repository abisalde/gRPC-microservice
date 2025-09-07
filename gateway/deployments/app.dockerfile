# Build stage

FROM golang:1.25-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Install git for go modules if private repos
RUN apk add --no-cache git gcc musl-dev bash

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy the entire project
COPY gateway gateway
COPY auth auth
COPY catalog catalog

# Copy air config from air.toml
COPY gateway/deployments/air.toml ./gateway/.air.toml



# Set working directory for Gateway service
WORKDIR /app/gateway

ENV AIR_WORKSPACE=/app/gateway \
    ENVIRONMENT=development

# Create directory for air builds
RUN mkdir -p /app/gateway/tmp

# Expose the port your server runs on
EXPOSE 8080


CMD ["air", "-c", "/app/gateway/.air.toml"]