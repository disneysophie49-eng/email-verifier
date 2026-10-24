
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /email-api ./cmd/apiserver

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /email-api /email-api

EXPOSE 8080

CMD ["/email-api"]
