# Step 1: Кэшируем зависимости
FROM golang:1.22-alpine3.19 AS modules
WORKDIR /modules
COPY go.mod go.sum /modules/
RUN go mod download

# Step 2: Сборка
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY --from=modules /go/pkg /go/pkg
COPY . /app 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd

# Step 3: Финальный образ
FROM alpine:3.19
WORKDIR /app
RUN ls -lah /app
COPY --from=builder /bin/app /app/
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN ls -lah /app/migrations

CMD ["/app/app"]
