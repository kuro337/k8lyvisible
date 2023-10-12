# Go Dockerfile

FROM golang:1.21-alpine3.18 as builder
WORKDIR /app
# Copy the go module and sum files to download dependencies
COPY src/go.mod src/go.sum /app/
RUN go mod download
# Copy the source files
COPY src/ /app/
# Build the application with optimizations
RUN go build -o /app/main .
# Use a smaller base image for the runtime
FROM golang:1.21-alpine3.18
# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main
# Expose the port the app runs on
EXPOSE 8080
# Specify the command to run on container start
CMD ["/app/main"]

