# Step 1: Use a base Go image
FROM golang:1.21-alpine

# Step 2: Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Step 3: Install necessary dependencies
RUN apk update && apk add --no-cache sqlite sqlite-libs gcc g++ musl-dev

# Step 4: Set the working directory inside the container
WORKDIR /app

# Step 5: Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 6: Copy application source code
COPY . .

# Step 7: Build the Go application
RUN go build -o main .

# Step 8: Expose the application port
EXPOSE 8080

# Step 9: Command to run the application
CMD ["./main"]
