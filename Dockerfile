# Stage 1: Builder
FROM golang:1.24.5-alpine AS builder
WORKDIR /app

RUN apk add --no-cache build-base git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -mod=mod -o /app/go-pharmacy-service ./cmd/main.go

# Stage 2: Runner
FROM alpine:latest

WORKDIR /app

RUN addgroup -g 1000 appuser && adduser -u 1000 -G appuser -D appuser

COPY --from=builder /app/go-pharmacy-service .

COPY .env .
COPY ./data ./data
RUN chown -R appuser:appuser /app/data
COPY ./docs ./docs
COPY ./data/pharmacies/pharmacies.csv /app/pharmacies.csv

EXPOSE 8080

USER appuser

ENTRYPOINT ["/app/go-pharmacy-service"]