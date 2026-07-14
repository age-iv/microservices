FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o censorship-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/censorship-service .
EXPOSE 8083
CMD ["./censorship-service"]
