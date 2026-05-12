FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o api_gateway ./cmd

FROM scratch

COPY --from=builder /app/api_gateway /

CMD ["/api_gateway"]