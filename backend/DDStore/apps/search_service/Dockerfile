FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o search_service ./cmd

FROM scratch

COPY --from=builder /app/search_service /

CMD ["/search_service"]
