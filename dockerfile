FROM golang:1.23-alpine

RUN apk update && apk add --no-cache sqlite sqlite-libs gcc g++ musl-dev

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

EXPOSE 8090

CMD ["./main"]

# to build : docker build -t test:3allal -f docker/dockerfile .
# to run   : docker run -d test:3allal
