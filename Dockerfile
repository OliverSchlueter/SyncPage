FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod file
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o syncpage .

# Verify the binary was created
RUN ls -l /app

# Use a minimal base image to keep the final image lightweight
FROM alpine:latest

# Set environment variable for the server port
ENV PORT=8080

# Create a directory for site files
RUN mkdir -p /app/data

# Copy the compiled binary from the builder stage
COPY --from=builder /app/syncpage /app/syncpage

# Verify the binary was copied
RUN ls -l /app

# Expose the application port
EXPOSE 8080

# Set the working directory
WORKDIR /app

# Run the application
CMD ["./syncpage"]
