FROM golang:alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o category-service cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/category-service .
CMD ["./category-service"]