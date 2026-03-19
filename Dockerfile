FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o health-aggregator main.go


FROM alpine:3.23

WORKDIR /root/

COPY --from=builder /app/health-aggregator .
COPY --from=builder /app/config.yaml .

CMD ["./health-aggregator"]