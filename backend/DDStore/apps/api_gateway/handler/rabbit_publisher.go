package handler

import (
	"context"
	"log"

	"github.com/Phuong-Hoang-Dai/DDStore/api_gateway/configs"
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



	RabbitHandler.conn = conn
	RabbitHandler.createdProductPublisher = publisher

	return conn
}

func (r RabbitPublisher) PublishProduct(message []byte) {
	if RabbitHandler.createdProductPublisher == nil || RabbitHandler.conn == nil {
		log.Panicf("Rabbit hasn't set up yet!")
	}

	res, err := RabbitHandler.createdProductPublisher.Publish(context.Background(), rmq.NewMessage(message))
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
