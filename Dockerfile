FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN apk add --no-cache build-base
COPY . .

RUN mkdir -p data/pharmacies
COPY ./data/pharmacies/pharmacies.csv ./data/pharmacies/pharmacies.csv

ENV CGO_ENABLED=1 
RUN go build -mod=mod -o /app/go-pharmacy-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-pharmacy-service .

COPY .env .
COPY ./data ./data
COPY ./docs ./docs

EXPOSE 8080

ENTRYPOINT ["/app/go-pharmacy-service"]