# Build stage
FROM golang:1.21-alpine AS builder

# Disable checksum verification globally BEFORE any go command
ENV GOSUMDB=off
ENV GOPROXY=direct

WORKDIR /app

# Install system dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go.mod only
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /hr-system ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies including curl for healthcheck
RUN apk --no-cache add ca-certificates tzdata curl

# Copy binary
COPY --from=builder /hr-system /app/hr-system

# Copy templates and static
COPY templates /app/templates
COPY static /app/static

# Create uploads dir
RUN mkdir -p /app/static/uploads

EXPOSE 8080

CMD ["/app/hr-system"]
