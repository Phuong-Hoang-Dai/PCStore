FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o payment_service ./cmd

FROM scratch

COPY --from=builder /app/payment_service /

CMD ["/payment_service"]