# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" \
    -o /bin/link-rpc ./services/link-rpc

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata && \
  addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /bin/link-rpc /app/link-rpc
COPY etc/link-rpc.compose.yaml /app/etc/link-rpc.compose.yaml

USER appuser

EXPOSE 9090 9101

ENTRYPOINT ["/app/link-rpc", "-f", "/app/etc/link-rpc.compose.yaml"]
