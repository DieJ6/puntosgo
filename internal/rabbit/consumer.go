package rabbit

import (
	"github.com/nmarsollier/commongo/log"
	"github.com/streadway/amqp"

	"github.com/DieJ6/puntosgo/internal/usecases"
)

// Estas constantes ya están en connection.go:
// const (
//     ExchangeName = "puntos_exchange"
//     QueueName    = "puntos_queue"
// )
// Acá solo dejamos la routing key específica de este consumer.

const (
	RkConsultaCompra = "consulta_compra"
)

type Consumer struct {
	conn             *amqp.Connection
	log              log.LogRusEntry
	ProcesarCompraUC *usecases.ProcesarCompraUC
}

func NewConsumer(
	conn *amqp.Connection,
	log log.LogRusEntry,
	procesarUC *usecases.ProcesarCompraUC,
) *Consumer {
	return &Consumer{
		conn:             conn,
		log:              log,
		ProcesarCompraUC: procesarUC,
	}
}

func (c *Consumer) Start() {
	ch, err := c.conn.Channel()
	if err != nil {
		c.log.Error(err)
		return
	}
	defer ch.Close()

	// Declarar el exchange (por si no está creado)
	if err := ch.ExchangeDeclare(
		ExchangeName, // viene de connection.go
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		c.log.Error(err)
		return
	}

	// Declarar la cola
	_, err = ch.QueueDeclare(
		QueueName, // viene de connection.go
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.log.Error(err)
		return
	}

	// Bind cola ↔ routing key
	if err := ch.QueueBind(
		QueueName,
		RkConsultaCompra,
		ExchangeName,
		false,
		nil,
	); err != nil {
		c.log.Error(err)
		return
	}

	msgs, err := ch.Consume(
		QueueName,
		"",
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.log.Error(err)
		return
	}

	c.log.Info("Rabbit consumer de puntosgo iniciado y escuchando consulta_compra")

	for msg := range msgs {
		switch msg.RoutingKey {
		case RkConsultaCompra:
			if err := c.ProcesarCompraUC.Consume(msg.Body); err != nil {
				c.log.Error(err)
				_ = msg.Nack(false, false)
				continue
			}
			_ = msg.Ack(false)
		default:
			_ = msg.Ack(false)
		}
	}
}
