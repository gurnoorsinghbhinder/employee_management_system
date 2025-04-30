FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o employee-management ./main.go

#final stage
FROM alpine:3.21
COPY --from=builder /app/employee-management /app/employee-management
COPY --from=builder /app/.env /app/.env
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set Docker environment flag
ENV DOCKER_ENVIRONMENT=true

#expose application port
EXPOSE 8080

WORKDIR /app

CMD ["./employee-management"]