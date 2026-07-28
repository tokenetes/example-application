# syntax=docker/dockerfile:1
# TODO: Use a multi-stage build process to reduce the final image size.
FROM --platform=$BUILDPLATFORM golang:1.25.0 AS builder

WORKDIR /app

# Install corporate CA cert if provided (e.g. Zscaler) without committing it to the repo.
# Build with: docker build --secret id=ca_cert,src=/path/to/corp-ca.pem ...
RUN --mount=type=secret,id=ca_cert,required=false \
    if [ -f /run/secrets/ca_cert ]; then \
      cp /run/secrets/ca_cert /usr/local/share/ca-certificates/corp-ca.crt && \
      update-ca-certificates; \
    fi

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd

CMD ["./service"]
