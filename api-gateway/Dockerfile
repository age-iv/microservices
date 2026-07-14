FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
