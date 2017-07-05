package amqprpc

import (
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
	"sync"
)

type pendingRequest struct {
	resultChan chan []byte
	errorChan  chan error
}

func newPendingRequest() *pendingRequest {
	return &pendingRequest{
		resultChan: make(chan []byte),
		errorChan:  make(chan error),
	}
}

func (p *pendingRequest) done() {
	close(p.resultChan)
	close(p.errorChan)
}

func (p *pendingRequest) res(r []byte) {
	p.resultChan <- r
	p.done()
}

func (p *pendingRequest) err(e error) {
	p.errorChan <- e
	p.done()
}

func NewAmqpClient(cfg *AMQPConfig, queue string) (*AmqpClient, error) {
	a := &AmqpClient{}
	err := a.init(cfg, queue)
	if err != nil {
		return nil, err
	}

	return a, nil
}

type AmqpClient struct {
	cfg       *AMQPConfig
	conn      *amqp.Connection
	ch        *amqp.Channel
	q         amqp.Queue
	msgs      <-chan amqp.Delivery
	sendLock  *sync.RWMutex
	pending   map[string]*pendingRequest
	workQueue string
}

func (ad *AmqpClient) init(cfg *AMQPConfig, queue string) error {
	ad.cfg = cfg
	ad.workQueue = queue
	ad.pending = make(map[string]*pendingRequest)
	ad.sendLock = &sync.RWMutex{}

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
		conn.Close()
		ch.Close()
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
		conn.Close()
		ch.Close()
		return err
	}

	go func() {
		for d := range ad.msgs {
			if c, ok := ad.pending[d.CorrelationId]; ok {
				c.res(d.Body)
			}
		}
	}()

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

	pending := newPendingRequest()
	ad.sendLock.Lock()
	ad.pending[corrId] = pending
	ad.sendLock.Unlock()

	go func() {
		err := ad.ch.Publish(
			"",                     // exchange
			ad.workQueue,           // routing key
			ad.cfg.MandatryPublish, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       ad.q.Name,
				Body:          request,
			},
		)

		if err != nil {
			pending.err(err)
		}
	}()

	select {
	case result := <-pending.resultChan:
		return result, nil
	case err := <-pending.errorChan:
		return nil, err
	}
}
