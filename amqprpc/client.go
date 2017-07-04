package amqprpc

import (
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
)

type AmqpClient struct {
	cfg       *AMQPConfig
	conn      *amqp.Connection
	ch        *amqp.Channel
	q         amqp.Queue
	msgs      <-chan amqp.Delivery
	workQueue string
}

func (ad *AmqpClient) Init(cfg *AMQPConfig, queue string) error {
	ad.cfg = cfg
	ad.workQueue = queue

	return nil
}

func (ad *AmqpClient) Close() error {
	if ad.conn == nil {
		return fmt.Errorf("AMQP Connection not open")
	}

	ad.conn.Close()
	ad.ch.Close()

	ad.conn = nil
	ad.ch = nil
	return nil
}

func (ad *AmqpClient) Dial() error {
	if ad.conn != nil {
		return nil
	}

	conn, err := amqp.Dial(ad.cfg.Dsn())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	ad.conn = conn
	ad.ch = ch
	ad.msgs = msgs
	ad.q = q

	return nil
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (ad *AmqpClient) Request(request []byte) ([]byte, error) {

	corrId := randomString(32)

	err := ad.ch.Publish(
		"",           // exchange
		ad.workQueue, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",

			CorrelationId: corrId,
			ReplyTo:       ad.q.Name,
			Body:          request,
		},
	)

	if err != nil {
		return nil, err
	}

	for d := range ad.msgs {
		if corrId == d.CorrelationId {
			return d.Body, nil
		}
	}

	return nil, fmt.Errorf("final message, didnt match")
}
