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
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/repos"
	"github.com/Phuong-Hoang-Dai/DDStore/app/product_service/internal/service"
	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RabbitHandler struct {
	productService service.ProductService
	conn           *rmq.AmqpConnection
}

func NewRabbitHandler(db *mongo.Client, redis *redis.Client, ctx context.Context) RabbitHandler {
	col := db.Database(configs.Cfg.DBName).Collection("products")
	repo := repos.NewMongoProductRepo(col)
	conn := SetupRabbitQueues(ctx)

	return RabbitHandler{productService: service.NewProductService(repo, redis), conn: conn}
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

	return conn
}

func (p RabbitHandler) ProductCreatedWorker(deliveries chan *rmq.IDeliveryContext, ctx context.Context) {
	for delivery := range deliveries {
		var data model.Product

		msg := (*delivery).Message()
		var body []byte
		if len(msg.Data) > 0 {
			body = msg.Data[0]
		}

		if err := json.Unmarshal(body, &data); err != nil {
			(*delivery).Discard(context.Background(), &amqp.Error{Condition: amqp.ErrCondDecodeError})
			return
		}

		_, err := p.productService.CreateProduct(ctx, data)

		if err != nil {
			(*delivery).Discard(ctx, &amqp.Error{Condition: amqp.ErrCondInternalError})
			return
		}

		if err := (*delivery).Accept(ctx); err != nil {
			log.Panicf("Failed to accept message: %v", err)
		}
	}
}

func (p RabbitHandler) ProductCreatedConsumer(queueName string, workerFunc workForQueue, ctx context.Context) {
	deliveries := make(chan *rmq.IDeliveryContext, configs.Cfg.WorkerPerQueue)
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
			log.Panicf("Failed to receive a message: %v", err)
		}
		deliveries <- &delivery
	}
}

func builddWorker(deliveries chan *rmq.IDeliveryContext, wg *sync.WaitGroup, workFunc workForQueue, ctx context.Context) {
	for range configs.Cfg.WorkerPerQueue {
		wg.Add(1)
		go workFunc(deliveries, ctx)
	}
}

type workForQueue func(msg chan *rmq.IDeliveryContext, ctx context.Context)
