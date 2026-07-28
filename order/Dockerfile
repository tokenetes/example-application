# TODO: Use a multi-stage build process to reduce the final image size.
FROM --platform=$BUILDPLATFORM golang:1.25.0 AS builder

WORKDIR /app

COPY go.* ./
COPY vendor/ vendor/

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -mod=vendor -o service ./cmd

CMD ["./service"]