package rabbit

import (
	"github.com/nmarsollier/commongo/log"
	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(routingKey string, body []byte) error
}

type publisher struct {
	conn *amqp.Connection
	log  log.LogRusEntry
}

func NewPublisher(conn *amqp.Connection, log log.LogRusEntry) Publisher {
	return &publisher{conn, log}
}

func (p *publisher) Publish(routingKey string, body []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		p.log.Error(err)
		return err
	}
	defer ch.Close()

	return ch.Publish(
		ExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
