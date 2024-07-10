package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"jpitservices.cloud/eda/internal"
)

func main() {

	conn, err := internal.ConnectRabbitMQ("percy", "secret", "localhost:5672", "customers")

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)

	if err != nil {
		panic(err)
	}

	defer client.Close()

	if err := client.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}
	if err := client.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}

	if err := client.CreateBiding("customers_created", "customers.created.*", "customer_events"); err != nil {
		panic(err)
	}

	if err := client.CreateBiding("customers_test", "customers.*", "customer_events"); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Send(ctx, "customer_events", "customers.created.se", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte("an cool message between services"),
	}); err != nil {
		panic(err)
	}

	if err := client.Send(ctx, "customer_events", "customers.test", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte("A second cool message"),
	}); err != nil {
		panic(err)
	}

	// time.Sleep(30 * time.Second)
	log.Println("client")

}
