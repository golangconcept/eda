package internal

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func ConnectRabbitMQ(username, pawword, host, vhost string) (*amqp.Connection, error) {

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, pawword, host, vhost))

	if err != nil {
		return nil, err
	}
	return conn, nil

}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {

	ch, err := conn.Channel()

	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, err
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)

	return err
}

func (rc RabbitClient) CreateBiding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	return rc.ch.PublishWithContext(ctx, exchange, routingKey, true, false, options)
}
