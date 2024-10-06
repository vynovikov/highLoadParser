# Stage 1: Build the Go application
FROM vynovikov/golang:0.1 AS builder

RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Set the working directory inside the container
WORKDIR /build

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download necessary Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -gcflags "all=-N -l" -o main cmd/highLoadParser/highLoadParser.go

# Final stage
FROM vynovikov/alpine:0.1

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /build/main /app/main
COPY --from=builder /go/bin/dlv /dlv
COPY --from=builder /build/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh
RUN chmod +x /dlv

# Expose the application port (optional, if it's an API)
EXPOSE 3000 40000

# Command to start Delve in headless mode for remote debugging
#CMD ["/dlv","--accept-multiclient","--continue", "--headless=true", "--listen=:40000", "--api-version=2","--check-go-version=false", "exec", "/app/main"]
CMD ["sh","./entrypoint.sh"]