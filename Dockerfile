# Use the official Alpine-based golang image as a parent image
FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /app

# Copy the current directory contents into the container
COPY . .

# Build the application
RUN go mod download && \
    CGO_ENABLED=0 go build -v -ldflags="-extldflags=-static" -tags netgo -a -o pixivfe

# Stage for creating the non-privileged user
FROM alpine:3.20 AS user-stage

RUN adduser -u 10001 -S pixivfe

# Stage for a smaller final image
FROM scratch

# Copy necessary files from the builder stage
COPY --from=builder /app/pixivfe /app/pixivfe
COPY --from=builder /app/assets /app/assets
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy passwd file for the non-privileged user from the user-stage
COPY --from=user-stage /etc/passwd /etc/passwd

# Set the working directory
WORKDIR /app

# Expose port 8282
EXPOSE 8282

# Switch to non-privileged user
USER pixivfe

# Set the entrypoint to the binary name
ENTRYPOINT ["/app/pixivfe"]
