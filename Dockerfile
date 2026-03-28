FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod ./

# Create go.sum with correct checksums
RUN go mod tidy

COPY . .

# Build from module root
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /hr-system .

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata curl

COPY --from=builder /hr-system /app/hr-system
COPY templates /app/templates
COPY static /app/static

RUN mkdir -p /app/static/uploads

EXPOSE 8080

CMD ["/app/hr-system"]
