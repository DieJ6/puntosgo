package rabbit

import (
	"github.com/nmarsollier/commongo/log"
	"github.com/streadway/amqp"

	"github.com/DieJ6/puntosgo/internal/usecases"
)

type publisher struct {
	conn *amqp.Connection
	log  log.LogRusEntry
}

// NewPublisher devuelve algo que implementa usecases.Publisher
func NewPublisher(conn *amqp.Connection, log log.LogRusEntry) usecases.Publisher {
	return &publisher{conn: conn, log: log}
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
