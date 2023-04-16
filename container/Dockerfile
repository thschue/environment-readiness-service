FROM golang:1.20.3-alpine3.16 AS builder

ENV CGO_ENABLED=0

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /workspace/environment-readiness-service ./

FROM gcr.io/distroless/static AS production

LABEL org.opencontainers.image.source="https://github.com/on-clouds/environment-readiness-service" \
    org.opencontainers.image.title="environment-readiness-service" \
    org.opencontainers.image.vendor="Thomas Schuetz" \
    org.opencontainers.image.licenses="MIT"

WORKDIR /
COPY --from=builder /workspace/environment-readiness-service .
USER 65532:65532

ENTRYPOINT ["/environment-readiness-service"]