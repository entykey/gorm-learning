# # Use latest Alpine-based Golang image
# ARG GO_VERSION=alpine
# FROM golang:${GO_VERSION} AS builder

# RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

# RUN mkdir -p /api
# WORKDIR /api

# # Copy module files first for better caching
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# # Copy the rest of the application files
# COPY . .

# # Build the application (fix: use correct filename)
# RUN go build -o ./app hello_server.go

# # Use smaller Alpine image as the final runtime image
# FROM alpine:latest

# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# RUN mkdir -p /api
# WORKDIR /api
# COPY --from=builder /api/app .

# EXPOSE 8080

# ENTRYPOINT ["./app"]



FROM golang:1.21-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Use the PORT environment variable
ENV PORT=3000

# Expose the internal port (3000)
EXPOSE 3000

ENTRYPOINT ["./main"]