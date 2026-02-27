FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o shortener ./cmd/app/main.go

FROM alpine:3.19 as runner

WORKDIR /app

RUN adduser -D shortener

COPY --from=builder /app/shortener .

USER shortener

EXPOSE 8888

CMD ["./shortener"]