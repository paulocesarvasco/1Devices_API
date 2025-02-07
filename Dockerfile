FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o devices-api ./cmd/server/main.go

FROM scratch

COPY  --from=builder /app/devices-api /devices-api

COPY --from=builder /app/static /static

EXPOSE 8080

CMD ["/devices-api"]
