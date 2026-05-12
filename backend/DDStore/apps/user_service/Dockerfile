FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o user_service ./cmd

FROM scratch

COPY --from=builder /app/user_service /

CMD ["/user_service"]