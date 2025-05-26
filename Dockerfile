# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o /daysleft

# Final stage
FROM alpine:latest

WORKDIR /

# Copy only the binary from builder
COPY --from=builder /daysleft /daysleft

# Run the binary
ENTRYPOINT ["/daysleft"]
