# Build context is ./apps (not ./apps/order_service): go.mod replaces the shared
# model module with ../model, so that directory must be inside the context.
FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY model ./model
COPY order_service/go.mod order_service/go.sum ./order_service/
WORKDIR /src/order_service
RUN go mod download

COPY order_service .
RUN CGO_ENABLED=0 go build -o /order_service ./cmd

FROM scratch

# CA certificates are required to call the Stripe API over HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /order_service /

CMD ["/order_service"]
