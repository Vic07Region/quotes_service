FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /quotes_service ./cmd/main.go

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /quotes_service /quotes_service

EXPOSE 8080

ENTRYPOINT ["/quotes_service"]