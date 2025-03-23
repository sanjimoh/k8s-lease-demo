FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o leader-election

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/leader-election .
ENTRYPOINT ["./leader-election"] 