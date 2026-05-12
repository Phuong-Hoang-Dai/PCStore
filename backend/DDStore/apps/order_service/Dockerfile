FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o order_service ./cmd

FROM scratch

COPY --from=builder /app/order_service /

CMD ["/order_service"]