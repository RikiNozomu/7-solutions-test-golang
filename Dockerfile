FROM golang:1.20-alpine AS builder

# install tools needed for fetching modules
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy rest and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o /app/app main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/app /app/app

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/app"]
