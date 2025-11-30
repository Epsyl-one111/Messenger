FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/web/static ./web/static/
COPY --from=builder /app/web/templates ./web/templates/
COPY --from=builder /app/.env . 
EXPOSE 8080
CMD ["./main"]