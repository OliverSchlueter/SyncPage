FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o syncpage .

# Use a minimal base image to keep the final image lightweight
FROM alpine:latest

# Set environment variable for the server port
ENV PORT=8080

# Create a directory for site files
RUN mkdir -p /app/sites

# Copy the compiled binary from the builder stage
COPY --from=builder /app/syncpage /app/syncpage

# Expose the application port
EXPOSE 8080

# Set the working directory
WORKDIR /app

# Run the application
CMD ["./syncpage"]
