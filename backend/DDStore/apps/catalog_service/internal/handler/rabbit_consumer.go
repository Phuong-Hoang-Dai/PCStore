package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/Azure/go-amqp"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/configs"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/model"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type RabbitConsumerModel struct {
	productService service.ProductService
	conn           *rmq.AmqpConnection
}

func NewRabbitConsumer(conn *rmq.AmqpConnection, productService service.ProductService) RabbitConsumerModel {
	return RabbitConsumerModel{productService: productService, conn: conn}
}

func SetupRabbitQueues(ctx context.Context) *rmq.AmqpConnection {
	env := rmq.NewEnvironment(configs.Cfg.BrokerURI, nil)
	conn, err := env.NewConnection(ctx)
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %v", err)
	}

	//declare queue
	if _, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: configs.Cfg.ProductCreatedQueue}); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}
	if _, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: configs.Cfg.ProductEditedQueue}); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}
	if _, err = conn.Management().DeclareQueue(ctx, &rmq.QuorumQueueSpecification{Name: configs.Cfg.ProductIndexQueue}); err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	return conn
}

func (p RabbitConsumerModel) ProductCreatedWorker(deliveries chan rmq.IDeliveryContext, ctx context.Context) {
	for delivery := range deliveries {
		var data model.Product

		msg := delivery.Message()
		var body []byte
		if len(msg.Data) > 0 {
			body = msg.Data[0]
		}

		// continue, not return: returning would stop this worker for good.
		if err := json.Unmarshal(body, &data); err != nil {
			delivery.Discard(context.Background(), &amqp.Error{Condition: amqp.ErrCondDecodeError})
			continue
		}

		id, err := p.productService.CreateProduct(ctx, data)

		if err != nil {
			delivery.Discard(ctx, &amqp.Error{Condition: amqp.ErrCondInternalError})
			continue
		}

		// CreateProduct also publishes the new product for search indexing.
		if err := delivery.Accept(ctx); err != nil {
			log.Printf("Failed to accept message: %v", err)
		}
		log.Printf("Created successfully: id = %s", id)
	}
}

func (p RabbitConsumerModel) CreateConsumer(queueName string, workerFunc workForQueue, ctx context.Context) {
	deliveries := make(chan rmq.IDeliveryContext, configs.Cfg.WorkerPerQueue)
	var wg sync.WaitGroup
	builddWorker(deliveries, &wg, workerFunc, ctx)

	consumer, err := p.conn.NewConsumer(ctx, queueName, nil)

	if err != nil {
		log.Panicf("Failed to create consumer: %v", err)
	}
	log.Printf("%s Waiting for messages. To exit press CTRL+C", queueName)
	for {
		delivery, err := consumer.Receive(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("Shutting down gracefully...")
				return
			}
			log.Printf("Failed to receive a message: %v", err)
			continue
		}
		deliveries <- delivery
	}
}

func builddWorker(deliveries chan rmq.IDeliveryContext, wg *sync.WaitGroup, workFunc workForQueue, ctx context.Context) {
	for range configs.Cfg.WorkerPerQueue {
		wg.Add(1)
		go workFunc(deliveries, ctx)
	}
}

type workForQueue func(msg chan rmq.IDeliveryContext, ctx context.Context)
