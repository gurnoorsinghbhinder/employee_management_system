FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env /app/.env

# Accepting build arguments
ARG MONGO_URI
ARG DATABASE_NAME
ARG COLLECTION_NAME

# Use build arguments
ENV MONGO_URI=${MONGO_URI}
ENV DATABASE_NAME=${DATABASE_NAME}
ENV COLLECTION_NAME=${COLLECTION_NAME}

# Build the Go app
RUN go build -o employee-management-app ./main.go

#final stage
FROM alpine:3.21
COPY --from=builder /app/employee-management-app /app/employee-management-app
COPY --from=builder /app/.env /app/.env
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

#set docker env flag
ENV DOCKER_ENV=true

#expose the port
EXPOSE 8080

WORKDIR /app

CMD ["./employee-management-app"]
