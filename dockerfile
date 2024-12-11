FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/main .
COPY --from=builder /lib /lib

EXPOSE 8089

CMD ["./main"]
