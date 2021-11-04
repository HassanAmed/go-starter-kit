# Start from the latest golang base image
FROM golang:1.17.2

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy Go Modules dependency requirements file
COPY go.mod .

# Copy Go Modules expected hashes file
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy all the app sources (recursively copies files and directories from the host into the image)
COPY . .

# Set http port
ENV PORT 4000

# Build the app
RUN go build

# Make port 5000 available to the world outside this container
EXPOSE $PORT

# Run the app
CMD ["./bin/go-starter-kit"]
