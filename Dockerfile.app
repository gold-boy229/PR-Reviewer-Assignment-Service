# Stage 1: Build the Go application
FROM golang:1.25.1 AS builder
WORKDIR /app
# Copy the entire project context into the builder container
COPY . .
COPY ./empty.env . 

# Move into the subdirectory where main.go resides to build the binary
WORKDIR /app/cmd/app

# Go mod download is usually done in the root WORKDIR
# It's better to manage dependencies at the root
RUN cd /app && go mod download
RUN CGO_ENABLED=0 go build -o /app/main .

# Stage 2: Run the application in a minimal image
FROM alpine:latest
WORKDIR /app
# Copy the built binary from the builder stage
COPY --from=builder /app/main .
# Set the entry point to run the binary
CMD ["./main"]
