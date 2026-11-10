# Multi-stage Dockerfile untuk Go Fiber Backend GaluraSwara
# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git & build tools
RUN apk add --no-cache git ca-certificates

# Cache dependensi Go
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode program
COPY . .

# Kompilasi binary statis untuk Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o galuraswara-backend .

# Stage 2: Minimal runtime image (scratch / alpine)
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

# Salin binary dari stage builder
COPY --from=builder /app/galuraswara-backend /app/galuraswara-backend

# Expose port aplikasi (Fly.io default internal port 8080)
EXPOSE 8080

ENV PORT=8080

CMD ["/app/galuraswara-backend"]
