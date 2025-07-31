# Start from a small Go base image
FROM golang:1.24.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically
RUN CGO_ENABLED=0 GOOS=linux go build -o google-checker

# Small final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/google-checker .

EXPOSE 8080
CMD ["./google-checker"]
