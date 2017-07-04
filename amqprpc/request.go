package amqprpc

import "github.com/streadway/amqp"

type AmqpRequest struct {
	ch *amqp.Channel
	d  amqp.Delivery
}

func (r *AmqpRequest) Body() []byte {
	return r.d.Body
}

func (r *AmqpRequest) Respond(response []byte) error {
	err := r.ch.Publish(
		"",          // exchange
		r.d.ReplyTo, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: r.d.CorrelationId,
			Body:          response,
		})

	if err != nil {
		return err
	}

	r.d.Ack(false)

	return nil
}
