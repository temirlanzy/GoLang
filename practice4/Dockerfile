FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o myapp ./cmd/api/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/myapp /app/myapp

EXPOSE 8080
ENTRYPOINT ["/app/myapp"]