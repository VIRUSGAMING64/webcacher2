FROM golang:1.26.2-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/webcacher2 ./webcacher2.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/webcacher2 /app/webcacher2
COPY webcacher.conf /app/webcacher.conf
COPY public /app/public

RUN mkdir -p /app/.cache

VOLUME ["/app/.cache", "/app/cache.json", "/app/queue.json", "/app/stats.json"]
EXPOSE 8092

ENTRYPOINT ["/app/webcacher2"]
