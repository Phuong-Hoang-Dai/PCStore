package handler

import (
	"context"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type RabbitPublisherModel struct {
	createdProductPublisher *rmq.Publisher
	conn                    *rmq.AmqpConnection
}

var RabbitPublisher RabbitPublisherModel

func SetupPublisher(handler RabbitConsumerModel, ctx context.Context) {
	publisher, err := handler.conn.NewPublisher(ctx, &rmq.QueueAddress{Queue: configs.Cfg.ProductCreatedQueue}, nil)
	if err != nil {
		log.Panicf("Failed to create publisher: %v", err)
	}
	RabbitPublisher.conn = handler.conn
	RabbitPublisher.createdProductPublisher = publisher
}

func (r RabbitPublisherModel) PublishProduct(message []byte) {
	if RabbitPublisher.createdProductPublisher == nil || RabbitPublisher.conn == nil {
		log.Panicf("Rabbit hasn't set up yet!")
	}

	res, err := RabbitPublisher.createdProductPublisher.Publish(context.Background(), rmq.NewMessage(message))
	if err != nil {
		log.Panicf("Failed to publish a message: %v", err)
	}
	switch res.Outcome.(type) {
	case *rmq.StateAccepted:
	default:
		log.Panicf("Unexpected publish outcome: %v", res.Outcome)
	}
	log.Printf(" [x] Sent %s\n", message)
}
