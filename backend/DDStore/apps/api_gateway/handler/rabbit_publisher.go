package handler

import (
	"context"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/logger"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type RabbitPublisher struct {
	createdProductPublisher *rmq.Publisher
	conn                    *rmq.AmqpConnection
}

var RabbitHandler RabbitPublisher

func SetupRabbit(ctx context.Context) *rmq.AmqpConnection {
	env := rmq.NewEnvironment(configs.Cfg.BrokerURI, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %v", err)
	}

	//declare queue
	if _, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: configs.Cfg.ProductCreatedQueue}); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	publisher, err := conn.NewPublisher(ctx, &rmq.QueueAddress{Queue: configs.Cfg.ProductCreatedQueue}, nil)
	if err != nil {
		log.Panicf("Failed to create publisher: %v", err)
	}

	RabbitHandler.conn = conn
	RabbitHandler.createdProductPublisher = publisher

	return conn
}

// PublishProduct sends the product to catalog_service with the request's transaction id.
func (r RabbitPublisher) PublishProduct(ctx context.Context, message []byte) {
	if RabbitHandler.createdProductPublisher == nil || RabbitHandler.conn == nil {
		log.Panicf("Rabbit hasn't set up yet!")
	}

	if err := logger.Publish(ctx, RabbitHandler.createdProductPublisher, configs.Cfg.ProductCreatedQueue, message); err != nil {
		log.Panicf("Failed to publish a message: %v", err)
	}
}
