package rabbit

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

func TryConnectToQueue(rabbitAddress string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var connErr error

	for {
		newConn, err := amqp.Dial(rabbitAddress)

		if err == nil {
			conn = newConn
			break
		}

		connErr = err
		log.Printf("Connection to rabbitmq failed :%s\n", err)
		time.Sleep(3 * time.Second)
	}

	return conn, connErr
}

func GetChannelToConsume(conn *amqp.Connection, exchangeName string) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to rabbitmq channel")
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to declare an exchange")
	}

	readQueue, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to declare a queue")
	}

	err = ch.QueueBind(
		readQueue.Name, // queue name
		"",             // routing key
		exchangeName,   // exchange
		false,
		nil)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		readQueue.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a consumer")
	}

	return msgs, nil
}
