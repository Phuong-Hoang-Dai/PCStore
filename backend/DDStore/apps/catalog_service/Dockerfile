FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o product_service ./cmd

FROM scratch

COPY --from=builder /app/product_service /

CMD ["/product_service"]