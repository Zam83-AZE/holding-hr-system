# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy only go.mod (not go.sum to avoid checksum issues)
COPY go.mod ./

# Disable checksum verification for this build
ENV GOSUMDB=off
ENV GOPROXY=direct

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /hr-system ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates, timezone data and curl (for healthcheck)
RUN apk --no-cache add ca-certificates tzdata curl

# Copy binary from builder
COPY --from=builder /hr-system /app/hr-system

# Copy templates and static files
COPY templates /app/templates
COPY static /app/static

# Create uploads directory
RUN mkdir -p /app/static/uploads

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/hr-system"]
