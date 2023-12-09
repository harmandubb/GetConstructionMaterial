FROM golang:1.21 AS builder
WORKDIR /app

RUN apt update

RUN apt install -y git ca-certificates \
    && update-ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the source code
COPY ./Server ./Server
COPY main.go .
COPY ./GCalls ./GCalls
COPY ./Auth2 ./Auth2
COPY ./API ./API

# # Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

# # Start a new stage from scratch for a smaller final image
FROM ubuntu:latest  

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /root/

# # Copy the binary from the builder stage
COPY --from=builder /app/myapp .

# COPY ./Server/cert.pem ./cert.pem
# COPY ./Server/key.pem ./key.pem

EXPOSE 8080

# Run the binary
CMD ["./myapp"]

# CMD ["go", "run", "main.go"]
