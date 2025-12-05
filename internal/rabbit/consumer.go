package rabbit

import (
	"encoding/json"
	"fmt"

	"github.com/nmarsollier/commongo/log"
	"github.com/streadway/amqp"

	"github.com/tuusuario/puntosgo/internal/usecases"
)

type Consumer struct {
	Conn             *amqp.Connection
	Log              log.LogRusEntry
	ProcesarCompraUC *usecases.ProcesarCompraUC
}

func NewConsumer(conn *amqp.Connection, log log.LogRusEntry, uc *usecases.ProcesarCompraUC) *Consumer {
	return &Consumer{conn, log, uc}
}

func (c *Consumer) Start() error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		QueueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {

			rk := msg.RoutingKey
			c.Log.Infof("Mensaje recibido: %s", rk)

			if rk == "consulta_compra" {

				if err := c.ProcesarCompraUC.Consume(msg); err != nil {
					c.Log.Error("Error procesando compra: ", err)
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
				continue
			}

			// otros routing keys podrían manejarse aquí
			msg.Ack(false)
		}
	}()

	fmt.Println("Consumer RabbitMQ iniciado.")
	return nil
}
