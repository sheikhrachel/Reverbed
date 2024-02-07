FROM public.ecr.aws/docker/library/golang:1.22 as builder
WORKDIR /build
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/main /app/
COPY --from=builder /build/swaggerui /app/swaggerui
COPY --from=builder /build/bin /app/bin
RUN . ./bin/activate-hermit
EXPOSE 8080 8125 8126
ENTRYPOINT ["./main"]
