# First "layer" of the image, only used to build
FROM public.ecr.aws/docker/library/golang:1.21 as builder
# Create a build dir and copy all the things
RUN mkdir /build
COPY . /build
# Set the Current Working Directory inside the container
WORKDIR /build
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Build the Go app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
# Use a really lightweight base environment layer for running
FROM gcr.io/distroless/base:debug
# Copy everything we need to run the service into the dir "app" of the container
COPY --from=builder /build/main /app/
COPY --from=builder /build/swaggerui /app/swaggerui
# Command to run the executable
WORKDIR /app
RUN apt-get update && apt-get install -y ffmpeg
RUN which ffmpeg
# Expose port 8080 and 8125 to the outside world
# router
EXPOSE 8080
# datadog agent
EXPOSE 8125
EXPOSE 8126
ENTRYPOINT ["./main"]
