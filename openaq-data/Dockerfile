FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openaq-data main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/openaq-data /app/
COPY --from=builder /app/internal /app/internal
COPY --from=builder /app/go.mod /app/
COPY --from=builder /app/go.sum /app/
ENTRYPOINT ["/app/openaq-data"]
