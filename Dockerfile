FROM golang:1.13 AS builder
# Copy local code to the container image.
WORKDIR /app
COPY . .
# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o crudapp

FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/crudapp /crudapp
CMD ["/crudapp", "-connectURI=192.168.99.100:27017", "-database=person", "-dbtype=mongo"]
EXPOSE 8000 8888
