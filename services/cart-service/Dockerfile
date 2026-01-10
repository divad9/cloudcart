FROM golang:1.21-alpine

WORKDIR /app

# Install git (needed for Go dependencies)
RUN apk add --no-cache git

# Copy go.mod first
COPY go.mod ./

# Download dependencies (this will create go.sum)
RUN go mod download

# Copy the rest of the code
COPY . .

# Tidy modules and build
RUN go mod tidy && go build -o main .

EXPOSE 5003

CMD ["./main"]
