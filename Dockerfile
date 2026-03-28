FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Copy everything first
COPY . .

# This creates go.sum with correct checksums
RUN go mod tidy

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /hr-system ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata curl

COPY --from=builder /hr-system /app/hr-system
COPY templates /app/templates
COPY static /app/static

RUN mkdir -p /app/static/uploads

EXPOSE 8080

CMD ["/app/hr-system"]
