package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ProductIndexEvent is the message search_service consumes from ProductIndexQueue.
// Keep it in sync with search_service/internal/handler/rabbit_consumer.go.
type ProductIndexEvent struct {
	Action  string         `json:"action"`
	Id      string         `json:"id"`
	Product *model.Product `json:"product,omitempty"`
}

const (
	ActionUpsert = "upsert"
	ActionDelete = "delete"
)

// RabbitPublisher implements service.ProductEvents.
type RabbitPublisher struct {
	publisher *rmq.Publisher
}

func NewRabbitPublisher(conn *rmq.AmqpConnection, ctx context.Context) *RabbitPublisher {
	// Must not be ProductCreatedQueue: this service consumes that queue, so publishing back
	// to it would re-create every product forever.
	publisher, err := conn.NewPublisher(ctx, &rmq.QueueAddress{Queue: configs.Cfg.ProductIndexQueue}, nil)
	if err != nil {
		log.Panicf("Failed to create publisher: %v", err)
	}
	return &RabbitPublisher{publisher: publisher}
}

func (r *RabbitPublisher) ProductUpserted(ctx context.Context, product model.Product) error {
	return r.publish(ctx, ProductIndexEvent{Action: ActionUpsert, Id: product.Id.Hex(), Product: &product})
}

func (r *RabbitPublisher) ProductDeleted(ctx context.Context, id bson.ObjectID) error {
	return r.publish(ctx, ProductIndexEvent{Action: ActionDelete, Id: id.Hex()})
}

func (r *RabbitPublisher) publish(ctx context.Context, event ProductIndexEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	res, err := r.publisher.Publish(ctx, rmq.NewMessage(message))
	if err != nil {
		return err
	}
	switch res.Outcome.(type) {
	case *rmq.StateAccepted:
	default:
		return fmt.Errorf("unexpected publish outcome: %v", res.Outcome)
	}
	log.Printf(" [x] Sent %s %s", event.Action, event.Id)
	return nil
}
