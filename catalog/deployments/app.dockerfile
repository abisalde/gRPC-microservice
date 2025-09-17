# Build stage

FROM golang:1.25-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Install git for go modules if private repos
RUN apk add --no-cache git gcc musl-dev bash curl unzip \
    && curl -L https://github.com/fullstorydev/grpcurl/releases/download/v1.9.1/grpcurl_1.9.1_linux_x86_64.tar.gz \
    | tar -xz -C /usr/local/bin

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download



# Copy the entire project
COPY catalog/ ./catalog/


# Copy air config from air.toml
COPY catalog/deployments/air.toml ./catalog/.air.toml



# Set working directory to catalog service
WORKDIR /app/catalog

ENV AIR_WORKSPACE=/app/catalog \
    ENVIRONMENT=development

# Create directory for air builds
RUN mkdir -p /app/catalog/tmp

# Expose the port your server runs on
EXPOSE 50052


CMD ["air", "-c", "/app/catalog/.air.toml"]