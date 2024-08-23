FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o server ./cmd/server

# Step 2: Run the Go binary in an Alpine image
FROM alpine:3.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# Enable the server to be run.
RUN chmod +x ./server

# Expose port 5000 to the outside world
EXPOSE 5000

# Command to run the executable
CMD ["./server"]