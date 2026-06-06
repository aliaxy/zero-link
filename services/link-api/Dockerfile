# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" \
    -o /bin/link-api ./services/link-api

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata && \
  addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /bin/link-api /app/link-api
COPY etc/link-api.compose.yaml /app/etc/link-api.compose.yaml

USER appuser

EXPOSE 8080 9100

ENTRYPOINT ["/app/link-api", "-f", "/app/etc/link-api.compose.yaml"]
