# Build stage
FROM golang:1.24-alpine3.20

WORKDIR /app

# Copy go mod and download deps
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]
