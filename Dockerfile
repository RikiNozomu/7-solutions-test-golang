FROM golang:1.25.3-alpine AS builder

# Enviroment variables
ENV GIN_MODE=release
ENV APP_ENV=production

# Copy your Go application source code
WORKDIR /app
COPY go.mod go.sum ./

# cache deps
RUN go mod download

# copy rest and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app .

# Run testing
RUN go test -v ./...

EXPOSE 8080

CMD ["/bin/app"]
