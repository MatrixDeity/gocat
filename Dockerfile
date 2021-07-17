FROM golang:1.14.3-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o ./gocat ./cmd

FROM alpine:latest
COPY --from=builder /build/gocat /app/
CMD ["/app/gocat"]
