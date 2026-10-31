package logger

// Only in the services that use RabbitMQ (api_gateway, catalog_service, search_service); keep
// the copies the same.

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// queueProperty is the AMQP application property that carries the transaction id.
const queueProperty = "transaction_id"

// Publish sends body to queue through publisher, with ctx's transaction id attached, and logs
// the message at level 1.
func Publish(ctx context.Context, publisher *rmq.Publisher, queue string, body []byte) error {
	start := time.Now()
	msg := rmq.NewMessage(body)
	if id := TransactionID(ctx); id != "" {
		msg.ApplicationProperties = map[string]any{queueProperty: id}
	}

	res, err := publisher.Publish(ctx, msg)
	if err == nil {
		if _, ok := res.Outcome.(*rmq.StateAccepted); !ok {
			err = fmt.Errorf("unexpected publish outcome: %v", res.Outcome)
		}
	}
	outcome := "published"
	if err != nil {
		outcome = "failed"
	}
	logMessage(ctx, "queue_publish", queue, body, outcome, err, start)
	return err
}

// Received starts the handling of a message taken from queue. It returns a context with the
// transaction id the publisher attached (a new one if there is none), and a func to call with
// the outcome (accepted, discarded, requeued) once the message is dealt with, which logs the
// message at level 1.
func Received(ctx context.Context, queue string, delivery rmq.IDeliveryContext) (context.Context, func(outcome string, err error)) {
	start := time.Now()
	msg := delivery.Message()
	id, _ := msg.ApplicationProperties[queueProperty].(string)
	if id == "" {
		id = NewTransactionID()
	}
	ctx = WithTransactionID(ctx, id)

	var body []byte
	if len(msg.Data) > 0 {
		body = msg.Data[0]
	}
	return ctx, func(outcome string, err error) {
		logMessage(ctx, "queue_consume", queue, body, outcome, err, start)
	}
}

func logMessage(ctx context.Context, kind, queue string, body []byte, outcome string, err error, start time.Time) {
	level := slog.LevelInfo
	attrs := []any{
		"kind", kind,
		"queue", queue,
		"outcome", outcome,
		"message_body", formatBody("", body),
		"duration_ms", milliseconds(time.Since(start)),
	}
	if err != nil {
		level = slog.LevelError
		attrs = append(attrs, "error", err.Error())
	}
	slog.Log(ctx, level, fmt.Sprintf("%s %s %s", kind, queue, outcome), attrs...)
}
