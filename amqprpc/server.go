package amqprpc

import (
	"fmt"
	"github.com/btccom/rpctransport-go/rpc"
	"github.com/streadway/amqp"
)

func NewAmqpServer(cfg *AMQPConfig, queue string) (*AmqpServer, error) {
	return &AmqpServer{
		cfg:       cfg,
		prefetch:  1,
		workQueue: queue,
	}, nil
}

type AmqpServer struct {
	cfg       *AMQPConfig
	workQueue string
	conn      *amqp.Connection
	ch        *amqp.Channel
	q         amqp.Queue
	prefetch  int
	msgs      <-chan amqp.Delivery
}

func (ad *AmqpServer) Consume() <-chan rpc.ServerRequest {
	requests := make(chan rpc.ServerRequest)
	go func() {
		for d := range ad.msgs {
			requests <- &AmqpRequest{ad.ch, d}
		}
	}()
	return requests
}

func (ad *AmqpServer) Dial() error {
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
		ad.workQueue, // name
		false,        // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return err
	}

	err = ch.Qos(
		ad.prefetch, // prefetch count
		0,           // prefetch size
		false,       // global
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
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

func (ad *AmqpServer) Close() error {
	if ad.conn == nil {
		return fmt.Errorf("AMQP Connection not open")
	}

	ad.conn.Close()
	ad.ch.Close()

	ad.conn = nil
	ad.ch = nil
	return nil
}

func (ad *AmqpServer) SetPrefetch(prefetch int) {
	ad.prefetch = prefetch
}
