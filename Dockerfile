# Use the official Go image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.20 as builder

# Copy local code to the container image.
WORKDIR /app

# # Fetch dependencies.
# # If the Go mod and sum files are not changed, then the Docker cache 
# # will skip the go mod download step. This ensures that we only re-download
# # dependencies when they have changed and makes builds faster.
COPY go.* ./
RUN go mod download

# Copy the rest of the code
COPY . ./

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./app

# # Use a Docker multi-stage build to create a lean production image.
# # This step starts another build stage with a much smaller image.
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
