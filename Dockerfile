FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hako cmd/main.go

FROM alpine:latest

EXPOSE 8080

COPY --from=builder /app/hako /opt/hako

CMD ["/opt/hako", "start", "-p", "8080", "--verbose"]
