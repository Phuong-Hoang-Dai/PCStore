package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/search_service/internal/service"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

// productIndexEvent is published by catalog_service on every product change.
// Keep it in sync with catalog_service/internal/handler/rabbit_publisher.go.
type productIndexEvent struct {
	Action  string         `json:"action"` // "upsert" or "delete"
	Id      string         `json:"id"`
	Product *model.Product `json:"product,omitempty"`
}

// retryDelay keeps a failing message (e.g. Elasticsearch down) from being redelivered
// in a hot loop. RabbitMQ drops it after the quorum queue's delivery limit.
const retryDelay = 2 * time.Second

type ProductIndexConsumer struct {
	service service.SearchService
}

func NewProductIndexConsumer(searchService service.SearchService) ProductIndexConsumer {
	return ProductIndexConsumer{service: searchService}
}

// Run consumes the index queue until ctx is cancelled. Messages are handled one at a
// time so an upsert and a delete of the same product are applied in order.
func (c ProductIndexConsumer) Run(ctx context.Context) error {
	env := rmq.NewEnvironment(configs.Cfg.BrokerURI, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		return fmt.Errorf("connect to RabbitMQ: %w", err)
	}
	defer conn.Close(context.Background())

	if _, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: configs.Cfg.ProductIndexQueue}); err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	consumer, err := conn.NewConsumer(ctx, configs.Cfg.ProductIndexQueue, nil)
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}
	log.Printf("%s Waiting for messages", configs.Cfg.ProductIndexQueue)

	for {
		delivery, err := consumer.Receive(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			log.Printf("Failed to receive a message: %v", err)
			continue
		}
		c.handle(ctx, delivery)
	}
}

func (c ProductIndexConsumer) handle(ctx context.Context, delivery rmq.IDeliveryContext) {
	var body []byte
	if data := delivery.Message().Data; len(data) > 0 {
		body = data[0]
	}

	var event productIndexEvent
	if err := json.Unmarshal(body, &event); err != nil || event.Id == "" {
		log.Printf("Discarding malformed index event: %s", body)
		delivery.Discard(ctx, &amqp.Error{Condition: amqp.ErrCondDecodeError})
		return
	}

	var err error
	switch {
	case event.Action == "upsert" && event.Product != nil:
		event.Product.Id = event.Id
		_, err = c.service.IndexProduct(ctx, *event.Product)
	case event.Action == "delete":
		err = c.service.DeleteProduct(ctx, event.Id)
	default:
		log.Printf("Discarding index event with unknown action %q", event.Action)
		delivery.Discard(ctx, &amqp.Error{Condition: amqp.ErrCondNotImplemented})
		return
	}

	if err != nil {
		log.Printf("Failed to %s product %s, retrying: %v", event.Action, event.Id, err)
		time.Sleep(retryDelay)
		delivery.Requeue(ctx)
		return
	}

	if err := delivery.Accept(ctx); err != nil {
		log.Printf("Failed to accept message: %v", err)
	}
	log.Printf("Indexed: %s %s", event.Action, event.Id)
}
